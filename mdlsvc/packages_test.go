package mdlsvc

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	gentypes "goa.design/model/mdlsvc/gen/types"
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

func Test_ListPackages(t *testing.T) {
	testCases := []struct {
		name     string
		dir      string
		expected []*gentypes.Package
		err      error
	}{
		{
			name: "Find Model Packages",
			dir:  "testdata",
			expected: []*gentypes.Package{
				{Dir: "testdata/model", ImportPath: "test/model/model"},
				{Dir: "testdata/parent/model", ImportPath: "test/parent/model/model"},
				{Dir: "testdata/parent/other_model", ImportPath: "test/parent/other_model/model"},
			},
			err: nil,
		},
		{
			name:     "No Package",
			dir:      "/does-not-exist-dir",
			expected: nil,
			err:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := &Service{dir: tc.dir}
			result, err := svc.ListPackages(context.Background())
			if result != nil {
				sort.Slice(result, func(i, j int) bool {
					return result[i].Dir < result[j].Dir
				})
			}
			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)
		})
	}
}
