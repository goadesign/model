package main

import (
	"fmt"

	"goa.design/model/codegen"
)

// Executes the DSL and uploads the corresponding workspace to Structurizr.
func main() {
	// Run the model DSL
	js, err := codegen.JSON("goa.design/model/examples/json/model", true)
	if err != nil {
		panic(err)
	}
	// Print the JSON
	fmt.Println(string(js))
}
