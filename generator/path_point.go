package generator

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// PathPoint represents one point in the path "chain".
type PathPoint struct {
	Name       string
	Level      int
	IsParam    bool
	Operations map[string]*openapi3.Operation // Method to OperationID
	Segments   map[string]*PathPoint
}

// Route returns the method name of the route
func (pp *PathPoint) Route() string {
	name := strings.Title(strings.Replace(strings.Replace(pp.Name, "{", "", -1), "}", "", -1))
	return fmt.Sprintf("routeLevel%d%s", pp.Level, name)
}

// AnySegmentIsParam returns tru if any of the segments is a path parameter.
func (pp *PathPoint) AnySegmentIsParam() bool {
	for _, pathPoint := range pp.Segments {
		if pathPoint.IsParam {
			return true
		}
	}
	return false
}
