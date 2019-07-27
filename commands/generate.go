package commands

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/michaelsauter/go-oas-server/generator"
)

// Generate renders Go files based on specification in file into directory outputDir.
func Generate(file string, outputDir string) error {
	spec, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile(file)
	if err != nil {
		return fmt.Errorf("could not load OpenAPI specification from %s: %s", file, err)
	}

	return generator.New(spec).Render(outputDir)
}
