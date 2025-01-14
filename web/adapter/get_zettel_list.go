//-----------------------------------------------------------------------------
// Copyright (c) 2020 Detlef Stern
//
// This file is part of zettelstore.
//
// Zettelstore is free software: you can redistribute it and/or modify it under
// the terms of the GNU Affero General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// Zettelstore is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with Zettelstore. If not, see <http://www.gnu.org/licenses/>.
//-----------------------------------------------------------------------------

// Package adapter provides handlers for web requests.
package adapter

import (
	"fmt"
	"log"
	"net/http"

	"zettelstore.de/z/ast"
	"zettelstore.de/z/config"
	"zettelstore.de/z/domain"
	"zettelstore.de/z/encoder"
	"zettelstore.de/z/parser"
	"zettelstore.de/z/usecase"
)

// MakeListMetaHandler creates a new HTTP handler for the use case "list some zettel".
func MakeListMetaHandler(te *TemplateEngine, listMeta usecase.ListMeta) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter, sorter := getFilterSorter(r)
		metaList, err := listMeta.Run(r.Context(), filter, sorter)
		if err != nil {
			checkUsecaseError(w, err)
			return
		}

		format := getFormat(r, encoder.GetDefaultFormat())
		w.Header().Set("Content-Type", format2ContentType(format))
		switch format {
		case "html":
			renderListMetaHTML(w, metaList)
		case "json", "djson":
			enc := encoder.Create(format)
			renderListMetaJSON(w, metaList, enc, format)
		case "native", "raw", "text", "zmk":
			http.Error(w, fmt.Sprintf("Zettel list in format %q not yet implemented", format), http.StatusNotImplemented)
			log.Println(format)
		default:
			http.Error(w, fmt.Sprintf("Zettel list not available in format %q", format), http.StatusBadRequest)
		}
	}
}

func renderListMetaHTML(w http.ResponseWriter, metaList []*domain.Meta) {
	buf := encoder.NewBufWriter(w)

	buf.WriteStrings("<html lang=\"", config.GetDefaultLang(), "\">\n<body>\n<ul>\n")
	for _, meta := range metaList {
		title := meta.GetDefault(domain.MetaKeyTitle, "")
		htmlTitle, err := formatInlines(parser.ParseTitle(title), "html")
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		buf.WriteStrings(
			"<li><a href=\"", urlForZettel('z', meta.Zid), "?_format=html", "\">",
			htmlTitle, "</a></li>\n")
	}
	buf.WriteString("</ul>\n</body>\n</html>")
	buf.Flush()
}

func renderListMetaJSON(w http.ResponseWriter, metaList []*domain.Meta, enc encoder.Encoder, format string) {
	if enc == nil {
		return
	}
	addFormat := format != encoder.GetDefaultFormat()
	detail := format == "djson"
	buf := encoder.NewBufWriter(w)
	buf.WriteString("{\"list\":[")
	for i, meta := range metaList {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteStrings("{\"id\":\"", meta.Zid.Format(), "\",\"url\":\"", urlForZettel('z', meta.Zid))
		if addFormat {
			buf.WriteStrings("?_format=", format)
		}
		buf.WriteString("\",\"meta\":")
		var title ast.InlineSlice
		if detail {
			title = parser.ParseTitle(meta.GetDefault(domain.MetaKeyTitle, ""))
		}
		enc.WriteMeta(&buf, meta, title)
		buf.WriteByte('}')
	}
	buf.WriteString("]}")
	buf.Flush()
}
