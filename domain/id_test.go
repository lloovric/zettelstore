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

// Package domain_test provides unit tests for testing domain specific functions.
package domain_test

import (
	"testing"

	"zettelstore.de/z/domain"
)

func TestParseZettelID(t *testing.T) {
}

func TestIsValid(t *testing.T) {
	validIDs := []string{
		"00000000000001",
		"00000000000020",
		"00000000000300",
		"00000000004000",
		"00000000050000",
		"00000000600000",
		"00000007000000",
		"00000080000000",
		"00000900000000",
		"00001000000000",
		"00020000000000",
		"00300000000000",
		"04000000000000",
		"50000000000000",
		"99999999999999",
		"00001007030200",
		"20200310195100",
	}

	for i, sid := range validIDs {
		zid, err := domain.ParseZettelID(sid)
		if err != nil {
			t.Errorf("i=%d: sid=%q is not valid, but should be. err=%v", i, sid, err)
		}
		s := zid.Format()
		if s != sid {
			t.Errorf("i=%d: zid=%v does not format to %q, but to %q", i, sid, zid.Format(), s)
		}
	}

	invalidIDs := []string{
		"", "0", "a",
		"00000000000000",
		"000000000000000",
		"99999999999999a",
		"20200310T195100",
	}

	for i, zid := range invalidIDs {
		if _, err := domain.ParseZettelID(zid); err == nil {
			t.Errorf("i=%d: zid=%q is valid, but should not be", i, zid)
		}
	}
}
