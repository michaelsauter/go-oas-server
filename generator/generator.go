package generator

import (
	"path"
	"strings"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
)

// Generator of specification to Go code
type Generator struct {
	spec *openapi3.Swagger
}

// New returns a new Generator
func New(spec *openapi3.Swagger) *Generator {
	return &Generator{spec: spec}
}

// Render renders specification into outputDir
func (g *Generator) Render(outputDir string) error {
	var err error

	packageName := path.Base(outputDir)

	pathHandlers, operations := preparePaths(g.spec)

	funcMap := template.FuncMap{
		"title":            strings.Title,
		"publicGoName":     publicGoName,
		"curled":           curled,
		"goTypeFrom":       goTypeFrom,
		"schemaType":       schemaType,
		"isPathParam":      isPathParam,
		"parametersByType": parametersByType,
		"extractParameter": extractParameter,
	}

	err = renderTemplateToFile(
		templateComponents,
		funcMap,
		struct {
			PackageName string
			Components  openapi3.Components
		}{
			PackageName: packageName,
			Components:  g.spec.Components,
		},
		outputDir+"/components.go",
	)
	if err != nil {
		return err
	}

	err = renderTemplateToFile(
		templateServer,
		funcMap,
		struct {
			PackageName string
			Operations  []*openapi3.Operation
			PathPoint   *PathPoint
		}{
			PackageName: packageName,
			Operations:  operations,
			PathPoint:   pathHandlers["routeLevel0Root"],
		},
		outputDir+"/server.go",
	)
	if err != nil {
		return err
	}

	err = renderTemplateToFile(
		templateRouting,
		funcMap,
		struct {
			PackageName  string
			PathHandlers map[string]*PathPoint
		}{
			PackageName:  packageName,
			PathHandlers: pathHandlers,
		},
		outputDir+"/routing.go",
	)
	if err != nil {
		return err
	}

	err = renderTemplateToFile(
		templateEndpoints,
		funcMap,
		struct {
			PackageName string
			Paths       openapi3.Paths
		}{
			PackageName: packageName,
			Paths:       g.spec.Paths,
		},
		outputDir+"/endpoints.go",
	)
	if err != nil {
		return err
	}

	return nil
}

func preparePaths(spec *openapi3.Swagger) (map[string]*PathPoint, []*openapi3.Operation) {
	pathHandlers := map[string]*PathPoint{}
	operations := []*openapi3.Operation{}

	pp := &PathPoint{
		Name:       "root",
		Level:      0,
		IsParam:    false,
		Operations: map[string]*openapi3.Operation{},
		Segments:   map[string]*PathPoint{},
	}

	pathHandlers[pp.Route()] = pp

	for k, v := range spec.Paths {
		parts := strings.Split(k, "/")
		maxLevel := len(parts) - 1

		ppp := pp
		for level, vp := range parts {
			if vp == "" {
				continue
			}

			if _, ok := ppp.Segments[vp]; ok {
				ppp = ppp.Segments[vp]
				if maxLevel == level {
					ppp.Operations = v.Operations()
					for _, op := range ppp.Operations {
						operations = append(operations, op)
					}
				}
			} else {
				ops := map[string]*openapi3.Operation{}
				if maxLevel == level {
					ops = v.Operations()
				}
				pathPoint := &PathPoint{
					Name:       vp,
					Level:      level,
					IsParam:    isPathParam(vp),
					Operations: ops,
					Segments:   map[string]*PathPoint{},
				}

				ppp.Segments[vp] = pathPoint

				pathHandlers[pathPoint.Route()] = pathPoint

				for _, op := range ops {
					operations = append(operations, op)
				}

				ppp = ppp.Segments[vp]
			}
		}
	}

	return pathHandlers, operations
}
