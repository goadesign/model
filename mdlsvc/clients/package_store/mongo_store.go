package pstore

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goa.design/clue/log"

	"goa.design/model/codegen"
	gentypes "goa.design/model/mdlsvc/gen/types"
)

type (
	// mongoStore provides an implementation of the package store client that uses MongoDB
	// as a backend.
	mongoStore struct {
		// Package store collection
		packages *mongo.Collection
		// directory used to compile models
		dir string
		// debug mode
		debug bool
	}

	// packageRow is the MongoDB representation of a model package.
	packageRow struct {
		Workspace string            `bson:"workspace"`
		Import    string            `bson:"import"`
		Dir       string            `bson:"dir"`
		CreatedAt time.Time         `bson:"createdAt"`
		UpdatedAt time.Time         `bson:"updatedAt"`
		Files     map[string]string `bson:"files"`
	}
)

// Database is the name of the MongoDB database used to store packages.
const Database = "packages"

// Collection is the name of the MongoDB collection used to store packages.
const Collection = "packages"

// Schema is the JSON schema used to validate packages.
var Schema = bson.M{
	"bsonType": "object",
	"properties": bson.M{
		"_id":       bson.M{"bsonType": "string"},
		"workspace": bson.M{"bsonType": "string"},
		"import":    bson.M{"bsonType": "string"},
		"dir":       bson.M{"bsonType": "string"},
		"createdAt": bson.M{"bsonType": "date"},
		"updatedAt": bson.M{"bsonType": "date"},
		"files": bson.M{
			"bsonType": "object",
			"patternProperties": bson.M{
				`.*\.go`: bson.M{
					"bsonType": "string",
				},
			},
		},
	},
	"required": bson.A{"workspace", "import", "dir", "createdAt", "updatedAt"},
}

// NewMongoClient returns a MongoDB client. compilationDir is the path to the
// directory used to compile models to render the corresponding JSON.
func NewMongoStore(ctx context.Context, mc *mongo.Client, compilationDir string, debug bool) (PackageStore, error) {
	packages, err := initPackagesCollection(ctx, mc)
	if err != nil {
		return nil, err
	}
	return &mongoStore{
		packages: packages,
		dir:      compilationDir,
		debug:    debug,
	}, nil
}

// ListWorkspaces returns the list of workspaces.
func (c *mongoStore) ListWorkspaces(ctx context.Context) ([]*gentypes.Workspace, error) {
	var ws []*gentypes.Workspace
	results, err := c.packages.Distinct(ctx, "workspace", bson.M{})
	if err != nil {
		log.Errorf(ctx, err, "MongoDB: failed to list workspaces")
		return nil, err
	}
	for _, res := range results {
		ws = append(ws, &gentypes.Workspace{
			Workspace: res.(string),
		})
	}
	return ws, nil
}

// CreatePackage creates a new package in the given workspace.
func (c *mongoStore) CreatePackage(ctx context.Context,  pkg *gentypes.PackageFile) error {
	// Create or update the package
	now := time.Now()
	_, err := c.packages.UpdateOne(
		ctx,
		bson.M{"workspace": pkg.Locator.Workspace, "dir": pkg.Locator.Dir},
		bson.M{
			"$setOnInsert": bson.M{
				"workspace": pkg.Locator.Workspace,
				"import":    pkg.
				"dir":       pkg.Dir,
				"createdAt": now,
				"updatedAt": now,
			},
			"$set": bson.M{
				"updatedAt": time.Now(),
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("MongoDB: failed to create or update package %s: %w", pkg.Dir, err)
	}
	return nil
}

// DeletePackage deletes the given package.
func (c *mongoStore) DeletePackage(ctx context.Context, loc *gentypes.PackageLocator) error {
	_, err := c.packages.DeleteOne(ctx, bson.M{"workspace": loc.Workspace, "dir": loc.Dir})
	if err != nil {
		log.Errorf(ctx, err, "MongoDB: failed to delete package %s", loc.Dir)
		return err
	}
	return nil
}

// ListPackages returns the list of model packages in the given workspace.
func (c *mongoStore) ListPackages(ctx context.Context, workspace string) ([]*gentypes.Package, error) {
	var pkgs []*gentypes.Package
	results, err := c.packages.Find(ctx, bson.M{"workspace": workspace})
	if err != nil {
		log.Errorf(ctx, err, "MongoDB: failed to list packages")
		return nil, err
	}
	defer results.Close(ctx)
	for results.Next(ctx) {
		var data bson.M
		if err := results.Decode(&data); err != nil {
			log.Errorf(ctx, err, "MongoDB: failed to decode package")
			return nil, err
		}
		pkgs = append(pkgs, &gentypes.Package{
			Dir:        data["dir"].(string),
			ImportPath: data["import"].(string),
		})
	}
	if err := results.Err(); err != nil {
		log.Errorf(ctx, err, "MongoDB: failed to list packages")
		return nil, err
	}
	return pkgs, nil
}

// ReadPackageFiles returns the list of files in the given model
// package.
func (c *mongoStore) ReadPackageFiles(ctx context.Context, loc *gentypes.PackageLocator) ([]*gentypes.PackageFile, error) {
	var files []*gentypes.PackageFile
	results, err := c.packages.Find(ctx, bson.M{"workspace": loc.Workspace, "dir": loc.Dir})
	if err != nil {
		log.Errorf(ctx, err, "MongoDB: failed to load %s", loc.Dir)
		return nil, err
	}
	defer results.Close(ctx)
	for results.Next(ctx) {
		var data bson.M
		if err := results.Decode(&data); err != nil {
			log.Errorf(ctx, err, "MongoDB: failed to decode %s", loc.Dir)
			return nil, err
		}
		for _, f := range data["files"].(bson.M) {
			files = append(files, &gentypes.PackageFile{
				Locator: &gentypes.FileLocator{
					Workspace: loc.Workspace,
					Dir:       loc.Dir,
					Filename:  f.(bson.M)["name"].(string),
				},
				Content: f.(bson.M)["content"].(string),
			})
		}
	}
	if err := results.Err(); err != nil {
		log.Errorf(ctx, err, "MongoDB: failed to load %s", loc.Dir)
		return nil, err
	}
	return files, nil
}

// Save persists the content of the given file for future
// retrieval with LoadPackageFiles.
func (c *mongoStore) Save(ctx context.Context, file *gentypes.PackageFile) error {
	_, err := c.packages.UpdateOne(
		ctx,
		bson.M{"workspace": file.Locator.Workspace, "dir": file.Locator.Dir},
		bson.M{
			"$set": bson.M{
				fmt.Sprintf("files.%s.content", file.Locator.Filename): file.Content,
			},
		},
	)
	if err != nil {
		log.Errorf(ctx, err, "MongoDB: failed to save %s", file.Locator.Filename)
		return err
	}
	return nil
}

// Subscribe notifies of updates to the model implemented by the given package.
// The returned channel streams JSON representations of the model. The
// subscription is closed when the context is canceled.
func (c *mongoStore) Subscribe(ctx context.Context, pkg *gentypes.PackageLocator) (chan []byte, error) {
	// Retrieve package path
	var pkgRow bson.M
	err := c.packages.FindOne(ctx, bson.M{"workspace": pkg.Workspace, "dir": pkg.Dir}).Decode(&pkgRow)
	if err != nil {
		return nil, fmt.Errorf("MongoDB: failed to subscribe to %s, failed to retrieve package import path: %w", pkg.Dir, err)
	}
	pkgPath := pkgRow["import"].(string)
	ch := make(chan []byte)
	match := bson.D{{
		"$match", bson.D{
			{"operationType", "update"},
			{"fullDocument.workspace", pkg.Workspace},
			{"fullDocument.dir", pkg.Dir},
		},
	}}
	changeStream, err := c.packages.Watch(ctx, mongo.Pipeline{match})
	if err != nil {
		return nil, fmt.Errorf("MongoDB: failed to subscribe to %s: %w", pkg.Dir, err)
	}
	go func() {
		defer changeStream.Close(ctx)
		for changeStream.Next(ctx) {
			js, err := c.compile(ctx, pkgPath, pkg.Workspace, pkg.Dir)
			if err != nil {
				log.Errorf(ctx, err, "MongoDB: failed to compile %s", pkg.Dir)
				continue
			}
			ch <- js
		}
	}()
	return ch, nil
}

// initPackagesCollection ensures that the packages collection exists and is
// properly indexed.
func initPackagesCollection(ctx context.Context, mc *mongo.Client) (*mongo.Collection, error) {
	db := mc.Database(Database)
	spec, err := db.ListCollectionSpecifications(ctx, bson.M{"name": Collection})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	// If collection doesn't exist create one
	if len(spec) == 0 {
		err := db.CreateCollection(ctx, Collection)
		if err != nil {
			return nil, err
		}
	}

	// Upsert schema validation rules.
	err = db.RunCommand(ctx, bson.D{
		{Key: "collMod", Value: Collection},
		{Key: "validator", Value: bson.M{"$jsonSchema": Schema}},
	}).Err()
	if err != nil {
		return nil, err
	}

	packages := db.Collection(Collection)
	// Create the index on the packages collection.
	_, err = packages.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "workspace", Value: 1},
			{Key: "dir", Value: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	return packages, nil
}

// compile compiles the given model package and returns the corresponding JSON.
func (c *mongoStore) compile(ctx context.Context, pkgPath, workspace, modelDir string) ([]byte, error) {
	var pkg packageRow
	err := c.packages.FindOne(ctx, bson.M{"workspace": workspace, "dir": modelDir}).Decode(&pkg)
	if err != nil {
		return nil, fmt.Errorf("MongoDB: failed to load %s in %s: %w", modelDir, workspace, err)
	}
	fullPath := filepath.Join(c.dir, workspace, modelDir)
	_, err = os.Stat(fullPath)
	if err == nil {
		// Remove the directory
		if err := os.RemoveAll(fullPath); err != nil {
			log.Errorf(ctx, err, "MongoDB: failed to delete directory %s", fullPath)
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("MongoDB: failed to stat directory %s: %w", fullPath, err)
	}
	if err := os.MkdirAll(filepath.Join(c.dir, workspace, modelDir), 0755); err != nil {
		log.Errorf(ctx, err, "MongoDB: failed to create directory %s", filepath.Join(c.dir, workspace, modelDir))
		return nil, err
	}
	for name, content := range pkg.Files {
		if err := os.WriteFile(filepath.Join(fullPath, name), []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("MongoDB: failed to write file %s: %w", filepath.Join(fullPath, name), err)
		}
	}
	// Compile the model
	return codegen.JSON(fullPath, pkgPath, c.debug)
}
