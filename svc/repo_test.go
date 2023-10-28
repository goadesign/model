package svc

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	mockrepo "goa.design/model/svc/clients/repo/mocks"
	gentypes "goa.design/model/svc/gen/types"
)

func Test_CreatePackage(t *testing.T) {
	repo := "testdata"
	dir := "new"
	content := "package new"
	filename := "model.go"
	testCases := []struct {
		name     string
		payload  *gentypes.PackageFile
		result   error
		expected string
		err      string
	}{
		{
			name: "Success",
			payload: &gentypes.PackageFile{
				Locator: &gentypes.FileLocator{
					Repository: repo,
					Dir:        dir,
					Filename:   filename,
				},
				Content: content,
			},
			expected: "testdata/model/new/design.go",
		},
		{
			name: "Create Package Error",
			payload: &gentypes.PackageFile{
				Locator: &gentypes.FileLocator{
					Repository: repo,
					Dir:        dir,
					Filename:   filename,
				},
				Content: content,
			},
			result:   fmt.Errorf("create package error"),
			expected: "",
			err:      "failed to create package new: create package error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := mockrepo.NewRepoHandler(t)
			handler.AddCreatePackage(func(_ context.Context, pf *gentypes.PackageFile) error {
				assert.Equal(t, repo, pf.Locator.Repository)
				assert.Equal(t, dir, pf.Locator.Dir)
				assert.Equal(t, filename, pf.Locator.Filename)
				assert.Equal(t, content, pf.Content)
				return tc.result
			})
			svc := &Service{dir: tc.name, handler: handler}
			err := svc.CreatePackage(context.Background(), tc.payload)
			if tc.err != "" {
				assert.EqualError(t, err, tc.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
