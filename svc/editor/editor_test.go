package editor

import (
	"go/format"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gentypes "goa.design/model/svc/gen/types"
)

func TestParser_UpsertElement(t *testing.T) {
	pkgdir := "model"
	defaultLocator := &gentypes.FileLocator{
		Dir:      pkgdir,
		Filename: DefaultModelFilename,
	}
	tests := []struct {
		name     string
		kind     ElementKind
		path     string
		existing map[string]string // existing code by filename
		code     string
		expected *gentypes.PackageFile
	}{
		{
			name: "Create new model file",
			kind: SoftwareSystemKind,
			path: "NewSoftwareSystem",
			code: addSystemCode,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+addSystemCode+endBrackets),
			},
		},
		{
			name:     "Add new system at end",
			kind:     SoftwareSystemKind,
			path:     "NewSoftwareSystem",
			existing: map[string]string{DefaultModelFilename: contentHeader + existingSystem + endBrackets},
			code:     addSystemCode,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+existingSystem+"\n"+addSystemCode+endBrackets),
			},
		},
		{
			name:     "Add new system before Views()",
			kind:     SoftwareSystemKind,
			path:     "NewSoftwareSystem",
			existing: map[string]string{DefaultModelFilename: contentHeader + existingSystem + viewsCode + endBrackets},
			code:     addSystemCode,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+existingSystem+"\n\t"+addSystemCode+viewsCode+endBrackets),
			},
		},
		{
			name:     "Add new container",
			kind:     ContainerKind,
			path:     "ExistingSystem/NewContainer",
			existing: map[string]string{DefaultModelFilename: contentHeader + existingSystem + endBrackets},
			code:     addContainerCode,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+addedContainer+endBrackets),
			},
		},
		{
			name:     "Add new component",
			kind:     ComponentKind,
			path:     "ExistingSystem/ExistingContainer/NewComponent",
			existing: map[string]string{DefaultModelFilename: contentHeader + existingSystem + endBrackets},
			code:     addComponentCode,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+addedComponent+endBrackets),
			},
		},
		{
			name:     "Update existing system",
			kind:     SoftwareSystemKind,
			path:     "ExistingSystem",
			existing: map[string]string{DefaultModelFilename: contentHeader + existingSystem + endBrackets},
			code:     editSystemCode,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+editSystemCode+endBrackets),
			},
		},
		{
			name:     "Update empty system",
			kind:     SoftwareSystemKind,
			path:     "ExistingSystem",
			existing: map[string]string{DefaultModelFilename: contentHeader + existingEmptySystem + endBrackets},
			code:     editSystemCode,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+editSystemCode+endBrackets),
			},
		},
		{
			name:     "Update existing system with relationship using empty system",
			kind:     SoftwareSystemKind,
			path:     "ExistingSystem",
			existing: map[string]string{DefaultModelFilename: contentHeader + existingSystemWithRel + endBrackets},
			code:     editSystemEmptyCode,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+editEmptySystemCodeWithRel+endBrackets),
			},
		},
		{
			name:     "Update existing system with relationship",
			kind:     SoftwareSystemKind,
			path:     "ExistingSystem",
			existing: map[string]string{DefaultModelFilename: contentHeader + existingSystemWithRel + endBrackets},
			code:     editSystemCode,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+editSystemCodeWithRel+endBrackets),
			},
		},
		{
			name:     "Update existing system with multiple relationships",
			kind:     SoftwareSystemKind,
			path:     "ExistingSystem",
			existing: map[string]string{DefaultModelFilename: contentHeader + existingSystemWithRels + endBrackets},
			code:     editSystemCode,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+editSystemCodeWithRels+endBrackets),
			},
		},
		{
			name:     "Update existing container",
			kind:     ContainerKind,
			path:     "ExistingSystem/ExistingContainer",
			existing: map[string]string{DefaultModelFilename: contentHeader + existingSystem + endBrackets},
			code:     editContainerCode,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+editedContainerCode+endBrackets),
			},
		},
		{
			name:     "Update existing container with relationship",
			kind:     ContainerKind,
			path:     "ExistingSystem/ExistingContainer",
			existing: map[string]string{DefaultModelFilename: contentHeader + existingSystemWithRels + endBrackets},
			code:     editContainerCode,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+editedContainerWithRelsCode+endBrackets),
			},
		},
		{
			name:     "Update existing component",
			kind:     ComponentKind,
			path:     "ExistingSystem/ExistingContainer/ExistingComponent",
			existing: map[string]string{DefaultModelFilename: contentHeader + existingSystem + endBrackets},
			code:     editComponentCode,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+editedComponentCode+endBrackets),
			},
		},
		{
			name:     "Update existing component with relationship",
			kind:     ComponentKind,
			path:     "ExistingSystem/ExistingContainer/ExistingComponent",
			existing: map[string]string{DefaultModelFilename: contentHeader + existingSystemWithRels + endBrackets},
			code:     editComponentCode,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+editedComponentCodeWithRels+endBrackets),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpdir, err := os.MkdirTemp(t.TempDir(), "model_parser_test")
			require.NoError(t, err)
			defer os.RemoveAll(tmpdir) // nolint: errcheck
			err = os.MkdirAll(filepath.Join(tmpdir, pkgdir), 0755)
			require.NoError(t, err)
			for filename, content := range tt.existing {
				err = os.WriteFile(filepath.Join(tmpdir, pkgdir, filename), []byte(content), 0644)
				require.NoError(t, err)
			}
			tt.expected.Locator.Repository = tmpdir
			p := NewEditor(tmpdir, pkgdir)
			res, err := p.UpsertElementByPath(tt.kind, tt.path, tt.code)
			require.NoError(t, err)
			assert.Equal(t, tt.expected.Locator, res.Locator, "locator")
			assert.Equal(t, tt.expected.Content, res.Content, "content")
		})
	}
}

func Test_UpsertRelationship(t *testing.T) {
	pkgdir := "model"
	defaultLocator := &gentypes.FileLocator{
		Dir:      pkgdir,
		Filename: DefaultModelFilename,
	}
	tests := []struct {
		name     string
		srcPath  string
		destPath string
		existing map[string]string // existing code by filename
		code     string
		expected *gentypes.PackageFile
	}{
		{
			name:     "Add relationship to system",
			srcPath:  "ExistingSystem",
			destPath: "AnotherSystem",
			existing: map[string]string{DefaultModelFilename: contentHeader + existingTwoSystems + endBrackets},
			code:     `Uses("AnotherSystem", "Test Relationship")`,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+existingTwoSystemsWithRel+endBrackets),
			},
		},
		{
			name:     "Add relationship to container",
			srcPath:  "ExistingSystem/ExistingContainer",
			destPath: "AnotherSystem",
			existing: map[string]string{DefaultModelFilename: contentHeader + existingTwoSystems + endBrackets},
			code:     `Uses("AnotherSystem", "Test Relationship")`,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+existingTwoSystemsWithContainerRel+endBrackets),
			},
		},
		{
			name:     "Update relationship to system",
			srcPath:  "ExistingSystem",
			destPath: "AnotherSystem",
			existing: map[string]string{DefaultModelFilename: contentHeader + existingTwoSystemsWithRel + endBrackets},
			code:     `Uses("AnotherSystem", "Edited Relationship")`,
			expected: &gentypes.PackageFile{
				Locator: defaultLocator,
				Content: formatted(t, contentHeader+editedTwoSystemsWithRel+endBrackets),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpdir, err := os.MkdirTemp(t.TempDir(), "model_parser_test")
			require.NoError(t, err)
			defer os.RemoveAll(tmpdir) // nolint: errcheck
			err = os.MkdirAll(filepath.Join(tmpdir, pkgdir), 0755)
			require.NoError(t, err)
			for filename, content := range tt.existing {
				err = os.WriteFile(filepath.Join(tmpdir, pkgdir, filename), []byte(content), 0644)
				require.NoError(t, err)
			}
			tt.expected.Locator.Repository = tmpdir
			p := NewEditor(tmpdir, pkgdir)
			res, err := p.UpsertRelationship(tt.srcPath, tt.destPath, tt.code)
			require.NoError(t, err)
			assert.Equal(t, tt.expected.Locator, res.Locator, "locator")
			assert.Equal(t, tt.expected.Content, res.Content, "content")
		})
	}
}

func formatted(t *testing.T, code string) string {
	t.Helper()
	bytes, err := format.Source([]byte(code))
	require.NoError(t, err, "failed to format code:\n%s", code)
	return string(bytes)
}

const (
	contentHeader = `package model

import . "goa.design/model/dsl"

var _ = Design(func() {
`

	endBrackets = `
})`

	existingSystem = `SoftwareSystem("ExistingSystem", func() {
		Tag("BeforeContainer")
		Container("ExistingContainer", func() {
			Tag("BeforeComponent")
			Component("ExistingComponent", func() {
				Tag("Component")
			})
			Tag("AfterComponent")
		})
		Tag("AfterContainer")
	})`

	existingTwoSystems = `SoftwareSystem("ExistingSystem", func() {
		Tag("BeforeContainer")
		Container("ExistingContainer", func() {
			Tag("BeforeComponent")
			Component("ExistingComponent", func() {
				Tag("Component")
			})
			Tag("AfterComponent")
		})
		Tag("AfterContainer")
	})
	SoftwareSystem("AnotherSystem", func() {
		Tag("BeforeContainer")
		Container("AnotherContainer", func() {
			Tag("BeforeComponent")
			Component("AnotherComponent", func() {
				Tag("Component")
			})
			Tag("AfterComponent")
		})
		Tag("AfterContainer")
	})`

	existingEmptySystem = `SoftwareSystem("ExistingSystem")`

	existingSystemWithRel = `SoftwareSystem("ExistingSystem", func() {
		Tag("BeforeContainer")
		Container("ExistingContainer", func() {
			Tag("BeforeComponent")
			Component("ExistingComponent", func() {
				Tag("Component")
			})
			Tag("AfterComponent")
		})
		Tag("AfterContainer")
		Uses("AnotherSystem", "Test Relationship", "Go and Goa", Synchronous, func() {
			Tag("Relationship")
		})
	})`

	existingTwoSystemsWithRel = `SoftwareSystem("ExistingSystem", func() {
		Tag("BeforeContainer")
		Container("ExistingContainer", func() {
			Tag("BeforeComponent")
			Component("ExistingComponent", func() {
				Tag("Component")
			})
			Tag("AfterComponent")
		})
		Tag("AfterContainer")
		Uses("AnotherSystem", "Test Relationship")
	})
	SoftwareSystem("AnotherSystem", func() {
		Tag("BeforeContainer")
		Container("AnotherContainer", func() {
			Tag("BeforeComponent")
			Component("AnotherComponent", func() {
				Tag("Component")
			})
			Tag("AfterComponent")
		})
		Tag("AfterContainer")
	})`

	existingTwoSystemsWithContainerRel = `SoftwareSystem("ExistingSystem", func() {
		Tag("BeforeContainer")
		Container("ExistingContainer", func() {
			Tag("BeforeComponent")
			Component("ExistingComponent", func() {
				Tag("Component")
			})
			Tag("AfterComponent")
			Uses("AnotherSystem", "Test Relationship")
		})
		Tag("AfterContainer")
	})
	SoftwareSystem("AnotherSystem", func() {
		Tag("BeforeContainer")
		Container("AnotherContainer", func() {
			Tag("BeforeComponent")
			Component("AnotherComponent", func() {
				Tag("Component")
			})
			Tag("AfterComponent")
		})
		Tag("AfterContainer")
	})`

	existingSystemWithRels = `SoftwareSystem("ExistingSystem", func() {
		Tag("BeforeContainer")
		Container("ExistingContainer", func() {
			Tag("BeforeComponent")
			Component("ExistingComponent", func() {
				Uses("AnotherSystem", "Ignored Relationship", "Go and Goa", Synchronous, func() {
					Tag("Should not be copied")
				})
				Tag("Component")
			})
			Uses("AnotherSystem", "Ignored Relationship", "Go and Goa", Synchronous, func() {
				Tag("Should not be copied")
			})
			Tag("AfterComponent")
		})
		Tag("AfterContainer")
		Uses("AnotherSystem", "Test Relationship", "Go and Goa", Synchronous, func() {
			Tag("Relationship")
		})
		Uses("YetAnotherSystem", "Test Relationship 2", "Go and Goa", Synchronous, func() {
			Tag("Relationship2")
		})
	})`

	addSystemCode = `SoftwareSystem("NewSystem", func() {
	Tag("NewSystem")
})`

	addContainerCode = `Container("NewContainer", func() {
	Tag("NewContainer")
})`

	addedContainer = `SoftwareSystem("ExistingSystem", func() {
	Tag("BeforeContainer")
	Container("ExistingContainer", func() {
		Tag("BeforeComponent")
		Component("ExistingComponent", func() {
			Tag("Component")
		})
		Tag("AfterComponent")
	})
	Tag("AfterContainer")
	Container("NewContainer", func() {
		Tag("NewContainer")
	})
})`

	addComponentCode = `Component("NewComponent", func() {
	Tag("NewComponent")
})`

	addedComponent = `SoftwareSystem("ExistingSystem", func() {
	Tag("BeforeContainer")
	Container("ExistingContainer", func() {
		Tag("BeforeComponent")
		Component("ExistingComponent", func() {
			Tag("Component")
		})
		Tag("AfterComponent")
		Component("NewComponent", func() {
			Tag("NewComponent")
		})
	})
	Tag("AfterContainer")
})`

	editSystemCode = `SoftwareSystem("ExistingSystem", func() {
	Tag("EditedSystem")
})`

	editEmptySystemCodeWithRel = `SoftwareSystem("ExistingSystem", func() {
	Uses("AnotherSystem", "Test Relationship", "Go and Goa", Synchronous, func() {
		Tag("Relationship")
	})
})`

	editSystemEmptyCode = `SoftwareSystem("ExistingSystem")`

	editSystemCodeWithRel = `SoftwareSystem("ExistingSystem", func() {
	Tag("EditedSystem")
	Uses("AnotherSystem", "Test Relationship", "Go and Goa", Synchronous, func() {
		Tag("Relationship")
	})
})`

	editSystemCodeWithRels = `SoftwareSystem("ExistingSystem", func() {
	Tag("EditedSystem")
	Uses("AnotherSystem", "Test Relationship", "Go and Goa", Synchronous, func() {
		Tag("Relationship")
	})
	Uses("YetAnotherSystem", "Test Relationship 2", "Go and Goa", Synchronous, func() {
		Tag("Relationship2")
	})
})`

	editContainerCode = `Container("ExistingContainer", func() {
	Tag("EditedContainer")
})`

	editedContainerCode = `SoftwareSystem("ExistingSystem", func() {
	Tag("BeforeContainer")
	Container("ExistingContainer", func() {
		Tag("EditedContainer")
	})
	Tag("AfterContainer")
})`

	editedContainerWithRelsCode = `SoftwareSystem("ExistingSystem", func() {
		Tag("BeforeContainer")
		Container("ExistingContainer", func() {
			Tag("EditedContainer")
			Uses("AnotherSystem", "Ignored Relationship", "Go and Goa", Synchronous, func() {
				Tag("Should not be copied")
			})
		})
		Tag("AfterContainer")
		Uses("AnotherSystem", "Test Relationship", "Go and Goa", Synchronous, func() {
			Tag("Relationship")
		})
		Uses("YetAnotherSystem", "Test Relationship 2", "Go and Goa", Synchronous, func() {
			Tag("Relationship2")
		})
	})`

	editComponentCode = `Component("ExistingComponent", func() {
	Tag("EditedComponent")
})`

	editedComponentCode = `SoftwareSystem("ExistingSystem", func() {
	Tag("BeforeContainer")
	Container("ExistingContainer", func() {
		Tag("BeforeComponent")
		Component("ExistingComponent", func() {
			Tag("EditedComponent")
		})
		Tag("AfterComponent")
	})
	Tag("AfterContainer")
})`

	editedComponentCodeWithRels = `SoftwareSystem("ExistingSystem", func() {
		Tag("BeforeContainer")
		Container("ExistingContainer", func() {
			Tag("BeforeComponent")
			Component("ExistingComponent", func() {
				Tag("EditedComponent")
				Uses("AnotherSystem", "Ignored Relationship", "Go and Goa", Synchronous, func() {
					Tag("Should not be copied")
				})
			})
			Uses("AnotherSystem", "Ignored Relationship", "Go and Goa", Synchronous, func() {
				Tag("Should not be copied")
			})
			Tag("AfterComponent")
		})
		Tag("AfterContainer")
		Uses("AnotherSystem", "Test Relationship", "Go and Goa", Synchronous, func() {
			Tag("Relationship")
		})
		Uses("YetAnotherSystem", "Test Relationship 2", "Go and Goa", Synchronous, func() {
			Tag("Relationship2")
		})
	})`

	editedTwoSystemsWithRel = `SoftwareSystem("ExistingSystem", func() {
		Tag("BeforeContainer")
		Container("ExistingContainer", func() {
			Tag("BeforeComponent")
			Component("ExistingComponent", func() {
				Tag("Component")
			})
			Tag("AfterComponent")
		})
		Tag("AfterContainer")
		Uses("AnotherSystem", "Edited Relationship")
	})
	SoftwareSystem("AnotherSystem", func() {
		Tag("BeforeContainer")
		Container("AnotherContainer", func() {
			Tag("BeforeComponent")
			Component("AnotherComponent", func() {
				Tag("Component")
			})
			Tag("AfterComponent")
		})
		Tag("AfterContainer")
	})`

	viewsCode = `
	Views(func() {
		SystemContextView("ExistingSystem", "ExistingSystemContext", func() {
			Tag("ExistingSystemContext")
		})
	})`
)
