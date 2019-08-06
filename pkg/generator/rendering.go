package generator

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
)

func curled(s string) string {
	return "{" + s + "}"
}

func renderTemplateToString(templateString string, funcMap template.FuncMap, data interface{}) (string, error) {
	b := bytes.NewBuffer([]byte{})

	tmpl, err := template.New("").Funcs(funcMap).Parse(templateString)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(b, data)
	if err != nil {
		return "", fmt.Errorf("could not render template: %s", err)
	}

	return b.String(), nil
}

func renderTemplateToFile(template string, funcMap template.FuncMap, data interface{}, outputFilename string) error {
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("could not create output file %s: %s", outputFilename, err)
	}
	defer outputFile.Close()
	err = renderTemplate(template, funcMap, data, outputFile)
	if err != nil {
		return err
	}
	err = exec.Command("goimports", "-w", outputFilename).Run()
	if err != nil {
		return err
	}
	return exec.Command("gofmt", "-w", outputFilename).Run()
}

func renderTemplate(rawTemplate string, funcMap template.FuncMap, data interface{}, output io.Writer) error {
	// parse
	tmpl, err := parseTemplate(rawTemplate, funcMap)
	if err != nil {
		return fmt.Errorf("could not parse template: %s", err)
	}
	// execute
	err = tmpl.Execute(output, data)
	if err != nil {
		return fmt.Errorf("could not render template: %s", err)
	}
	return nil
}

func parseTemplate(rawTemplate string, funcMap template.FuncMap) (*template.Template, error) {
	tmpl, err := template.New("").Funcs(funcMap).Parse(rawTemplate)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func parametersByType(parameters openapi3.Parameters, t string) openapi3.Parameters {
	filteredParameters := openapi3.Parameters{}
	for _, p := range parameters {
		if p.Value.In == t {
			filteredParameters = append(filteredParameters, p)
		}
	}
	return filteredParameters
}

func isPathParam(s string) bool {
	return strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")
}

func publicGoName(s string) string {
	return strings.Title(strings.Replace(s, "-", "_", -1))
}

func schemaType(name string, prefix string, schema *openapi3.SchemaRef) string {
	funcMap := template.FuncMap{
		"title":        strings.Title,
		"publicGoName": publicGoName,
		"curled": func(s string) string {
			return "{" + s + "}"
		},
		"schemaRef": func(s string) string {
			sl := strings.Split(s, "/")
			return sl[len(sl)-1]
		},
		"goTypeFrom": goTypeFrom,
	}

	out, err := renderTemplateToString(
		templateType,
		funcMap,
		struct {
			Name   string
			Prefix string
			Schema *openapi3.SchemaRef
		}{
			Name:   name,
			Prefix: prefix,
			Schema: schema,
		})
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func extractParameter(location string, rawVariable string, param *openapi3.ParameterRef) string {
	funcMap := template.FuncMap{
		"title":        strings.Title,
		"publicGoName": publicGoName,
		"curled": func(s string) string {
			return "{" + s + "}"
		},
		"schemaRef": func(s string) string {
			sl := strings.Split(s, "/")
			return sl[len(sl)-1]
		},
		"goTypeFrom": goTypeFrom,
	}

	out, err := renderTemplateToString(
		templateParameterExtraction,
		funcMap,
		struct {
			Location    string
			RawVariable string
			Param       *openapi3.ParameterRef
		}{
			Location:    location,
			RawVariable: rawVariable,
			Param:       param,
		})
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func goTypeFrom(s *openapi3.SchemaRef) string {
	if s.Ref != "" {
		sl := strings.Split(s.Ref, "/")
		return "Model" + sl[len(sl)-1]
	}

	val := s.Value
	switch val.Type {
	case "number":
		return "float64"
	case "integer":
		return "int"
	case "boolean":
		return "bool"
	case "array":
		itemType := goTypeFrom(val.Items)
		return "[]" + itemType
	case "object":
		return "map[string]interface{}"
	case "string":
		switch val.Format {
		case "uuid":
			return "uuid.UUID"
		case "date-time":
			return "time.Time"
		case "date":
			return "time.Time"
		default:
			return "string"
		}
	default:
		return "interface{}"
	}
}
