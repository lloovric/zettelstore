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

// Package encoder provides a generic interface to encode the abstract syntax
// tree into some text form.
package encoder

import (
	"zettelstore.de/z/ast"
	"zettelstore.de/z/domain"
)

// StringOption is an option with a string value
type StringOption struct {
	Key   string
	Value string
}

// Name returns the visible name of this option.
func (so *StringOption) Name() string { return so.Key }

// BoolOption is an option with a boolean value.
type BoolOption struct {
	Key   string
	Value bool
}

// Name returns the visible name of this option.
func (bo *BoolOption) Name() string { return bo.Key }

// MetaOption is an option with meta data as the value.
type MetaOption struct {
	Meta *domain.Meta
}

// Name returns the visible name of this option.
func (mo *MetaOption) Name() string { return "meta" }

// StringsOption is an option that have a sequence of strings as the value.
type StringsOption struct {
	Key   string
	Value []string
}

// Name returns the visible name of this option.
func (so *StringsOption) Name() string { return so.Key }

// AdaptLinkOption specifies a link adapter.
type AdaptLinkOption struct {
	Adapter func(*ast.LinkNode) ast.InlineNode
}

// Name returns the visible name of this option.
func (al *AdaptLinkOption) Name() string { return "AdaptLinkOption" }

// AdaptImageOption specifies an image adapter.
type AdaptImageOption struct {
	Adapter func(*ast.ImageNode) ast.InlineNode
}

// Name returns the visible name of this option.
func (al *AdaptImageOption) Name() string { return "AdaptImageOption" }

// AdaptCiteOption specifies a citation adapter.
type AdaptCiteOption struct {
	Adapter func(*ast.CiteNode) ast.InlineNode
}

// Name returns the visible name of this option.
func (al *AdaptCiteOption) Name() string { return "AdaptCiteOption" }
