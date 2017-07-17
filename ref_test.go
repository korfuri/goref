package goref_test

import (
	"encoding/json"
	"testing"

	"github.com/korfuri/goref"
	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	correct := "{\"from\":{\"position\":{\"filename\":\"x/foo.go\",\"start_line\":42,\"end_line\":42,\"end_col\":10},\"package\":\"path/to/x\",\"ident\":\"foo\"},\"to\":{\"position\":{\"filename\":\"y/bar.go\",\"start_line\":314,\"start_col\":11,\"end_line\":314,\"end_col\":14},\"package\":\"path/to/y\",\"ident\":\"bar\"}}"
	r := goref.Ref{
		RefType: goref.Instantiation,
		FromPosition: goref.Position{
			File: "x/foo.go",
			PosC: 0,
			PosL: 42,
			EndC: 10,
			EndL: 42,
		},
		ToPosition: goref.Position{
			File: "y/bar.go",
			PosC: 11,
			PosL: 314,
			EndC: 14,
			EndL: 314,
		},
		FromIdent: "foo",
		ToIdent:   "bar",
		FromPackage: &goref.Package{
			Path: "path/to/x",
		},
		ToPackage: &goref.Package{
			Path: "path/to/y",
		},
	}
	j, err := json.Marshal(r)
	assert.NoError(t, err)
	assert.Equal(t, correct, string(j))
}
