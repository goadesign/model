# Big Bank PLC Example

This examples reproduces the Structurizr [Big Bank PLC](https://structurizr.com/share/36141) example with Model.
See
[model.go](https://github.com/goadesign/model/blob/master/examples/big_bank_plc/model/model.go)
for the complete DSL.

## Running

The example can be uploaded to the Structurizr service using the `stz`
command line tool (see the
[README](https://github.com/goadesign/model/tree/master/README.md) for
details on installing and using the tool).

```bash
stz gen goa.design/model/examples/big_bank_plc/model
```

This generates the file `model.json` that contains the JSON representation of
the design. This file can be uploaded to a Structurizr workspace:

```bash
stz put model.json -id XXX -key YYY -secret ZZZ
```
