## Editor webapp

The mdl editor is a web application that runs in the browser.
It is embedded in the go executable, so it can be served
as static files. `mdl` runs an HTTP server that serves the
editor as static files. `mdl` also serves the model and layout data
dynamically, for the editor to load. The editor then renders the
selected view as an SVG, allowing the user to edit positions of elements
and shapes of relationships.

### Development setup

To develop the mdl and the editor, start the TypeScript compiler in watch mode and run mdl go program in devmode.

`mdl` can be instructed to serve the editor files from disk instead
of the embedded copies, to allow for easy development.
```
DEVMODE=1 go run ./cmd/mdl ... mdl params
``` 

Compile and run the TypeScript application in watch mode
```
yarn install
yarn watch
```

`yarn watch` will watch for changes in the webapp files and recompile. 
Simply refresh the browser to see the changes.
