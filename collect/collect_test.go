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

// Package collect_test provides some unit test for collectors.
package collect_test

import (
	"testing"

	"zettelstore.de/z/ast"
	"zettelstore.de/z/collect"
)

func parseRef(s string) *ast.Reference {
	r := ast.ParseReference(s)
	if !r.IsValid() {
		panic(s)
	}
	return r
}

func TestLinks(t *testing.T) {
	zettel := &ast.Zettel{}
	links, images := collect.References(zettel)
	if links != nil || images != nil {
		t.Error("No links/images expected, but got:", links, "and", images)
	}

	intNode := &ast.LinkNode{Ref: parseRef("01234567890123")}
	para := &ast.ParaNode{
		Inlines: ast.InlineSlice{
			intNode,
			&ast.LinkNode{Ref: parseRef("https://zettelstore.de/z")},
		},
	}
	zettel.Ast = ast.BlockSlice{para}
	links, images = collect.References(zettel)
	if links == nil || images != nil {
		t.Error("Links expected, and no images, but got:", links, "and", images)
	}

	para.Inlines = append(para.Inlines, intNode)
	links, images = collect.References(zettel)
	if cnt := len(links); cnt != 3 {
		t.Error("Link count does not work. Expected: 3, got", links)
	}
}

func TestImage(t *testing.T) {
	zettel := &ast.Zettel{
		Ast: ast.BlockSlice{
			&ast.ParaNode{
				Inlines: ast.InlineSlice{
					&ast.ImageNode{Ref: parseRef("12345678901234")},
				},
			},
		},
	}
	_, images := collect.References(zettel)
	if images == nil {
		t.Error("Only image expected, but got: ", images)
	}
}
