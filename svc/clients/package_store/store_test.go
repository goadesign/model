package pstore

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_modulePath(t *testing.T) {
	testCases := []struct {
		name               string
		dir                string
		expectedFilePath   string
		expectedModulePath string
		err                error
	}{
		{
			"Valid Go Module",
			"testdata/model",
			"testdata/model/go.mod",
			"test/model",
			nil,
		},
		{
			"Valid Go Module",
			"testdata/parent/model",
			"testdata/parent/go.mod",
			"test/parent",
			nil,
		},
		{
			"Invalid Go Module",
			"testdata/invalid_module",
			"",
			"",
			fmt.Errorf("invalid go.mod file testdata/invalid_module/go.mod"),
		},
		{
			"No go.mod file",
			"/",
			"",
			"",
			fmt.Errorf("failed to find go.mod file"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fpath, mpath, err := modulePath(tc.dir)
			assert.Equal(t, tc.expectedFilePath, fpath)
			assert.Equal(t, tc.expectedModulePath, mpath)
			assert.Equal(t, tc.err, err)
		})
	}
}
