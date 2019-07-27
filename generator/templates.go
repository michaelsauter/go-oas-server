package generator

const (
	templateComponents = `// Code generated by go generate; DO NOT EDIT.
package {{ .PackageName }}

// Schemas
{{- range $csName, $csDefinition := .Components.Schemas }}
  {{- schemaType $csName "Model" $csDefinition }}
{{ end }}

// Responses
{{- range $resName, $resDefinition := .Components.Responses }}
  {{- $mediaTypeContent := $resDefinition.Value.Content.Get "application/json" }}
  {{- if $mediaTypeContent }}
    {{- schemaType $resName "Response" $mediaTypeContent.Schema }}
  {{- end }}
{{- end }}
`
	templateEndpoints = `// Code generated by go generate; DO NOT EDIT.
package {{ .PackageName }}

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

{{- range $pathName, $pathDefinition := .Paths }}

{{- $operations := $pathDefinition.Operations }}
{{- range $method, $op := $operations }}
// APIEndpoint{{ title .OperationID }} describes the {{ $op.OperationID }} endpoint.
type APIEndpoint{{ title $op.OperationID }} struct {
	handler     APIOperation{{title $op.OperationID}}
	middlewares Middlewares
}

{{- if $op.RequestBody }}
  {{- $mediaTypeContent :=  $op.RequestBody.Value.Content.Get "application/json"}}
  {{- if $mediaTypeContent }}
    {{- $rbName := title $op.OperationID }}
    {{- schemaType $rbName "RequestBody" $mediaTypeContent.Schema }}
  {{- else }}
  // Unsupported media type for request body, only application/json is supported
  {{- end }}
{{- end }}
{{- if $op.Parameters }}
// Parameters{{ title $op.OperationID }} describes the parameters for the {{.OperationID}} endpoint.
type Parameters{{ title $op.OperationID }} struct {
	{{- range $op.Parameters }}
	{{ publicGoName .Value.Name}} {{ goTypeFrom .Value.Schema }}
	{{- end }}
}
{{- end }}

// Boot{{ title $op.OperationID }} boots the {{ $op.OperationID }} endpoint.
func (oas *OpenAPIServer) Boot{{ title $op.OperationID }}() {
	m, h := oas.Server.Handle{{ title $op.OperationID }}()
	oas.{{ $op.OperationID }} = APIEndpoint{{ title $op.OperationID }}{
		handler:     h,
		middlewares: m,
	}
}

// APIOperation{{ title $op.OperationID }} is an alias for the func signature.
{{- if and $op.RequestBody $op.Parameters }}
type APIOperation{{ title $op.OperationID }} func(w http.ResponseWriter, r *http.Request, p Parameters{{ title $op.OperationID }}, rb RequestBody{{ title $op.OperationID }})
{{- else if $op.RequestBody }}
type APIOperation{{ title $op.OperationID }} func(w http.ResponseWriter, r *http.Request, rb RequestBody{{ title $op.OperationID }})
{{- else if $op.Parameters }}
type APIOperation{{ title $op.OperationID }} func(w http.ResponseWriter, r *http.Request, p Parameters{{ title $op.OperationID }})
{{- else }}
type APIOperation{{ title $op.OperationID }} func(w http.ResponseWriter, r *http.Request)
{{- end }}

func (e APIEndpoint{{ title $op.OperationID }}) execute(w http.ResponseWriter, r *http.Request) {
	{{- if $op.RequestBody }}
	// Request Body
	var rb RequestBody{{ title $op.OperationID }}
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		http.Error(w, "Could not decode JSON", http.StatusBadRequest)
		return
	}
	{{- end }}

	{{- if $op.Parameters }}
	p := Parameters{{ title $op.OperationID }}{}
	{{- end }}
	{{- $queryParameters := parametersByType $op.Parameters "query" }}
	{{- if $queryParameters }}
	// Query parameters
	var rawQueryValue string
	{{- range $queryParameters }}
	// Handle {{.Value.Name}} param
	rawQueryValue = r.URL.Query().Get("{{ .Value.Name }}")
	{{- extractParameter "Query" "rawQueryValue" . }}
	{{- end }}
	{{- end }}

	{{- $pathParameters := parametersByType .Parameters "path" }}
	{{- if $pathParameters }}
	// Path parameters
	{{- range $pathParameters }}
	oasPathParts := strings.Split("/pets/{{curled .Value.Name}}", "/")
	pathParts := strings.Split(r.URL.Path, "/")
	for k, v := range pathParts {
		if oasPathParts[k] == "{{curled .Value.Name}}" {
			{{- $goType := goTypeFrom .Value.Schema }}
			{{- if eq $goType "int" }}
			rawValue, err := strconv.Atoi(v)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid user id %q", v), http.StatusBadRequest)
				return
			}
			{{- else }}
			panic("Not implemented")
			{{- end }}
			p.{{ title .Value.Name }} = rawValue
		}
	}
	{{- end }}
	{{- end }}

	{{- $headerParameters := parametersByType .Parameters "header" }}
	{{- if $headerParameters }}
	// Header Parameters
	var rawHeaderValue string
	{{- range $headerParameters }}
	rawHeaderValue = r.Header.Get("{{.Value.Name}}")
	{{- extractParameter "Header" "rawHeaderValue" . }}
	{{- end }}
	{{- end }}

	// Call handler
	{{- if and $op.RequestBody $op.Parameters }}
	e.handler(w, r, p, rb)
	{{- else if $op.RequestBody }}
	e.handler(w, r, rb)
	{{- else if $op.Parameters }}
	e.handler(w, r, p)
	{{- else }}
	e.handler(w, r)
	{{- end }}
}
{{- end }}
{{- end }}
`
	templateRouting = `// Code generated by go generate; DO NOT EDIT.
package {{ .PackageName }}

import (
	"net/http"
	"path"
	"strings"
)

{{- range $route, $routePathPoint := .PathHandlers }}
func (oas *OpenAPIServer) {{ $route }}(w http.ResponseWriter, r *http.Request, p string) {
	var head string
	head, p = shiftPath(p)

	// If head is empty, try one of the methods
	if head == "" {
		{{- range $method, $operation := $routePathPoint.Operations}}
		if r.Method == "{{ $method }}" {
			oas.serve(oas.{{ $operation.OperationID }}.execute, oas.{{ $operation.OperationID }}.middlewares, w, r)
			return
		}
		{{- end }}
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	{{- if $routePathPoint.Segments }}
	// Match against segments
	{{- range $segment, $pathPoint := $routePathPoint.Segments }}
	{{- if not $pathPoint.IsParam }}
	if head == "{{ $segment }}" {
		oas.{{ $pathPoint.Route }}(w, r, p)
		return
	}
	{{- end }}
	{{- end }}
	{{- end }}

	{{- if $routePathPoint.AnySegmentIsParam }}
	// "Match" against path param
	{{- range $segment, $pathPoint := $routePathPoint.Segments}}
	{{- if $pathPoint.IsParam }}
	oas.{{ $pathPoint.Route }}(w, r, p)
	return
	{{- end }}
	{{- end }}
	{{- else }}
	// Not found
	http.Error(w, "Not Found", http.StatusNotFound)
	{{- end }}
}
{{- end }}

// shiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
// Use e.g. like this: head, req.URL.Path = ShiftPath(req.URL.Path)
// Taken from https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html.
func shiftPath(p string) (head, tail string) {
    p = path.Clean("/" + p)
    i := strings.Index(p[1:], "/") + 1
    if i <= 0 {
        return p[1:], "/"
    }
    return p[1:i], p[i:]
}
`
	templateServer = `// Code generated by go generate; DO NOT EDIT.
package {{ .PackageName }}

import (
	"net/http"
)

// OpenAPIServer wraps the custom server.
type OpenAPIServer struct {
	Server     CustomServer
	{{- range .Operations }}
	{{ .OperationID }}  APIEndpoint{{ title .OperationID }}
	{{- end }}
	middlewares Middlewares
}

// Middleware.
type Middleware func(next http.HandlerFunc) http.HandlerFunc

// Ordered collection of middlewares.
type Middlewares []Middleware

// CustomServer forces that all API operations and middlewares are implemented on the custom server.
type CustomServer interface {
	{{- range .Operations }}
	Handle{{title .OperationID}}() (Middlewares, APIOperation{{title .OperationID}})
	{{- end }}
	Middlewares() Middlewares
}

// NewOpenAPIServer returns an OpenAPIServer with initialised middlewares.
func NewOpenAPIServer(s CustomServer) *OpenAPIServer {
	return &OpenAPIServer{Server: s, middlewares: s.Middlewares()}
}

// ServeHTTP implements the http.Handler interface.
func (oas *OpenAPIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routingHandler := func (w http.ResponseWriter, r *http.Request) {
		oas.routeLevel0Root(w, r, r.URL.Path)
	}
	// Global middleware
	for i := range oas.middlewares {
		routingHandler = oas.middlewares[len(oas.middlewares)-1-i](routingHandler)
	}
	routingHandler(w, r)
}

// Boot boots the server.
func (oas *OpenAPIServer) Boot() *OpenAPIServer {
	{{- range .Operations }}
	oas.Boot{{title .OperationID}}()
	{{- end }}
	return oas
}

// Serve applies middleware for the endpoint.
func (oas *OpenAPIServer) serve(h http.HandlerFunc, m Middlewares, w http.ResponseWriter, r *http.Request) {
	for i := range m {
		h = m[len(m)-1-i](h)
	}
	h(w, r)
}

// And adds all middlewares to the end of the current middlewares.
func (m Middlewares) And(h ...Middleware) Middlewares {
	m = append(m, h...)
	return m
}
`
	templateType = `{{- if .Schema.Ref }}
type {{ .Prefix }}{{ title .Name }} {{ goTypeFrom .Schema }}
{{- else }}
{{- if eq .Schema.Value.Type "object" }}
{{- if .Schema.Value.Description }}
// {{ .Schema.Value.Description }}
{{- end }}
type {{ .Prefix }}{{ title .Name }} struct {
    {{- range $propName, $propDefinition := .Schema.Value.Properties }}
        {{ publicGoName $propName }} {{ goTypeFrom $propDefinition }} ` + "`" + `json:"{{ $propName }}"` + "`" + ` {{ if $propDefinition.Value.Description }}// {{ $propDefinition.Value.Description }}{{ end }}
    {{- end }}
}
{{- else }}
// Not supported yet!
{{- end }}
{{- end }}
`
	templateParameterExtraction = `{{- if .Param.Value.Required }}
if len({{ .RawVariable }}) == 0 {
    http.Error(w, "{{ .Location }} '{{ .Param.Value.Name }}' required", http.StatusBadRequest)
    return
}
{{- else if .Param.Value.Schema.Value.Default }}
p.{{ publicGoName .Param.Value.Name }} = {{ .Param.Value.Schema.Value.Default }}
{{- end }}
if len({{ .RawVariable }}) > 0 {
    {{- if eq .Param.Value.Schema.Value.Type "string" }}
        {{- if eq .Param.Value.Schema.Value.Format "uuid" }}
        u, err := uuid.Parse({{ .RawVariable }})
        if err != nil {
            http.Error(w, "{{ .Location }} '{{ .Param.Value.Name }}' is not a UUID", http.StatusBadRequest)
            return
        }
        p.{{ publicGoName .Param.Value.Name }} = u
        {{- else if eq .Param.Value.Schema.Value.Format "date-time" }}
        t, err := time.Parse(time.RFC3339, {{ .RawVariable }})
        if err != nil {
            http.Error(w, "{{ .Location }} '{{ .Param.Value.Name }}' is not a date-time", http.StatusBadRequest)
            return
        }
        p.{{ publicGoName .Param.Value.Name }} = t
        {{- else if eq .Param.Value.Schema.Value.Format "date" }}
        t, err := time.Parse("2006-01-02", {{ .RawVariable }})
        if err != nil {
            http.Error(w, "{{ .Location }} '{{ .Param.Value.Name }}' is not a date", http.StatusBadRequest)
            return
        }
        p.{{ publicGoName .Param.Value.Name }} = t
        {{- else }}
        {{- if .Param.Value.Schema.Value.MinLength }}
        if len({{ .RawVariable }}) < {{ .Param.Value.Schema.Value.MinLength }} {
            http.Error(w, "{{ .Location }} '{{ .Param.Value.Name }}' must be {{ .Param.Value.Schema.Value.MinLength }} characters long", http.StatusBadRequest)
            return
        }
        {{- end }}
        {{- if .Param.Value.Schema.Value.MaxLength }}
        if len({{ .RawVariable }}) > {{ .Param.Value.Schema.Value.MaxLength }} {
            http.Error(w, "{{ .Location }} '{{ .Param.Value.Name }}' must not be longer than {{ .Param.Value.Schema.Value.MaxLength }} characters", http.StatusBadRequest)
            return
        }
        {{- end }}
        p.{{ publicGoName .Param.Value.Name }} = {{ .RawVariable }}
        {{- end }}
    {{- else if eq .Param.Value.Schema.Value.Type "integer" }}
    intValue, err := strconv.Atoi({{ .RawVariable }})
    if err != nil {
        http.Error(w, "{{ .Location }} '{{ .Param.Value.Name }}' is not an integer", http.StatusBadRequest)
        return
    }
    {{- if .Param.Value.Schema.Value.Min }}
    if len({{ .RawVariable }}) < {{ .Param.Value.Schema.Value.Min }} {
        http.Error(w, "{{ .Location }} '{{ .Param.Value.Name }}' must be at least {{ .Param.Value.Schema.Value.Min }}", http.StatusBadRequest)
        return
    }
    {{- end }}
    {{- if .Param.Value.Schema.Value.Max }}
    if len({{ .RawVariable }}) > {{ .Param.Value.Schema.Value.Max }} {
        http.Error(w, "{{ .Location }} '{{ .Param.Value.Name }}' must be at most {{ .Param.Value.Schema.Value.Max }}", http.StatusBadRequest)
        return
    }
    {{- end }}
    p.{{ publicGoName .Param.Value.Name }} = intValue
    {{- end}}
}
`
)
