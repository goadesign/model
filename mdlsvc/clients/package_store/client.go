package pstore

import (
	"context"

	gentypes "goa.design/model/mdlsvc/gen/types"
)

type (
	// Client is the package store client used by the mdlsvc service.
	Client interface {
		// ListPackages returns the list of model packages in the given
		// workspace.
		ListPackages(context.Context, string) ([]*gentypes.Package, error)
		// ListPackageFiles returns the list of files in the given model
		// package.
		LoadPackageFiles(context.Context, *gentypes.PackageLocator) ([]gentypes.PackageFile, error)
		// Save persists the content of the given file for future
		// retrieval with LoadPackageFiles.
		Save(context.Context, *gentypes.PackageFile) error
		// Subscribe notifies of updates to the model implemented by the
		// given package. The returned channel streams JSON
		// representations of the model. Call the returned stop function
		// to stop listening.
		Subscribe(context.Context, *gentypes.PackageLocator) (chan []byte, func(), error)
	}
)
