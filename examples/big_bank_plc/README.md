# Big Bank PLC Example

This examples reproduces the Structurizr [Big Bank PLC]() example with Model.
See
[model.go](https://github.com/goadesign/model/blob/master/examples/big_bank_plc/model/model.go)
for the complete DSL.

## Running

Using the `stz` command line:

```bash
stz gen goa.design/model/examples/big_bank_plc/model
```

This generates the file `model.json` that contains the JSON representation of
the workspace.  This file can be used to upload the workspace, again using stz:

```bash
stz put model.json -id XXX -key YYY -secret ZZZ
```
