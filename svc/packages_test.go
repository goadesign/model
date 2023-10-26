package svc

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	gentypes "goa.design/model/svc/gen/types"
)

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
			result, err := svc.ListPackages(context.Background(), &gentypes.Workspace{Workspace: "testdata"})
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
