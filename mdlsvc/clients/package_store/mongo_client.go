package pstore

import (
	"context"
	"fmt"
	"os"
	"path"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goa.design/clue/log"

	"goa.design/model/codegen"
	gentypes "goa.design/model/mdlsvc/gen/types"
)

type (
	// mongoClient provides an implementation of the package store client that uses MongoDB
	// as a backend.
	mongoClient struct {
		// MongoDB client
		db *mongo.Client
		// Package store collection
		coll *mongo.Collection
		// directory used to compile models
		dir string
		// debug mode
		debug bool
	}
)

// PackageDatabase is the name of the MongoDB database used to store packages.
const PackageDatabase = "packages"

// PackageCollection is the name of the MongoDB collection used to store packages.
const PackageCollection = "packages"

// NewMongoClient returns a MongoDB client. compilationDir is the path to the
// directory used to compile models to render the corresponding JSON.
func NewMongoClient(db *mongo.Client, compilationDir string, debug bool) Client {
	coll := db.Database(PackageDatabase).Collection(PackageCollection)
	return &mongoClient{
		db:    db,
		coll:  coll,
		dir:   compilationDir,
		debug: debug,
	}
}

// Load implements Load.
func (c *mongoClient) Load(ctx context.Context, root string, p *gentypes.DSLFile) (string, error) {
	var res string
	err := c.coll.FindOne(ctx, bson.M{"root": root, "dir": p.Dir, "name": p.Name}).Decode(&res)
	if err != nil {
		log.Errorf(ctx, err, "MongoDB: failed to load %s", path.Join(p.Dir, p.Name))
		return "", err
	}
	return res, nil
}

// Save implements Save.
func (c *mongoClient) Save(ctx context.Context, root string, p *gentypes.ModelDSL) error {
	opts := options.Update().SetUpsert(true)
	doc := bson.M{"root": root, "dir": p.File.Dir, "name": p.File.Name, "content": p.Content}
	_, err := c.coll.UpdateOne(ctx, bson.M{"root": root, "dir": p.File.Dir, "name": p.File.Name}, bson.M{"$set": doc}, opts)
	if err != nil {
		log.Errorf(ctx, err, "MongoDB: failed to save %s", path.Join(p.File.Dir, p.File.Name))
		return err
	}
	return nil
}

// Subscribe implements Subscribe.
func (c *mongoClient) Subscribe(ctx context.Context, root string, p *gentypes.ModelDir, pkgPath string) (chan []byte, func(), error) {
	ch := make(chan []byte)
	stop := func() {
		close(ch)
	}
	match := bson.D{
		{
			"$match", bson.D{
				{"$or", bson.D{{"operationType", "update"}, {"operationType", "insert"}}},
				{"fullDocument.root", root},
				{"fullDocument.dir", p.Dir},
			},
		},
	}
	changeStream, err := c.coll.Watch(ctx, mongo.Pipeline{match})
	if err != nil {
		log.Errorf(ctx, err, "MongoDB: failed to subscribe to %s", p.Dir)
		return nil, nil, err
	}
	go func() {
		defer changeStream.Close(ctx)
		for changeStream.Next(ctx) {
			js, err := c.compile(ctx, pkgPath, root, p.Dir)
			if err != nil {
				log.Errorf(ctx, err, "MongoDB: failed to compile %s", p.Dir)
				continue
			}
			ch <- js
		}
	}()
	return ch, stop, nil
}

// compile compiles the given model package and returns the corresponding JSON.
func (c *mongoClient) compile(ctx context.Context, pkgPath, root, modelDir string) ([]byte, error) {
	files, err := c.coll.Find(ctx, bson.M{"root": root, "dir": modelDir})
	if err != nil {
		log.Errorf(ctx, err, "MongoDB: failed to load %s", modelDir)
		return nil, err
	}
	defer files.Close(ctx)
	var dslFiles map[string]string
	for files.Next(ctx) {
		var data bson.M
		if err := files.Decode(&data); err != nil {
			log.Errorf(ctx, err, "MongoDB: failed to decode %s", modelDir)
			return nil, err
		}
		dslFiles[data["name"].(string)] = data["content"].(string)
	}
	if err := files.Err(); err != nil {
		log.Errorf(ctx, err, "MongoDB: failed to load %s", modelDir)
		return nil, err
	}
	// Write the files to disk with a default go.mod
	tmpDir, err := writeFiles(c.dir, pkgPath, dslFiles)
	if err != nil {
		log.Errorf(ctx, err, "MongoDB: failed to write files to disk")
		return nil, err
	}
	defer func() {
		os.RemoveAll(tmpDir)
	}()
	// Compile the model
	return codegen.JSON(tmpDir, pkgPath, c.debug)
}

// writeFiles writes the given files to disk in a temporary directory under the
// given directory and returns the directory path.
func writeFiles(dir, pkgPath string, files map[string]string) (string, error) {
	tmpDir, err := os.MkdirTemp(dir, codegen.TmpDirPrefix)
	if err != nil {
		return "", err
	}
	for name, content := range files {
		if err := os.WriteFile(path.Join(tmpDir, name), []byte(content), 0644); err != nil {
			return "", err
		}
	}
	if err := os.WriteFile(path.Join(tmpDir, "go.mod"), []byte(fmt.Sprintf("module %s\n", pkgPath)), 0644); err != nil {
		return "", err
	}
	return tmpDir, nil
}
