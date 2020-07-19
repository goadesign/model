# Big Bank PLC Example

This examples reproduces the Structurizr [Big Bank PLC]() example with
Structurizr for Go.  See
[model.go](https://github.com/goadesign/structurizr/blob/master/examples/big_bank_plc/model/model.go)
for the complete DSL.

## Running

Using the `stz` command line:

```bash
stz gen goa.design/structurizr/examples/big_bank_plc/model
```

This generates the file `model.json` that contains the JSON representation of
the workspace.  This file can be used to upload the workspace, again using stz:

```bash
stz put model.json -workspace XXX -key YYY -secret ZZZ
```
