# go-oas-server

Generate Go server code from an [OpenAPI 3 specification](https://swagger.io/specification/).

This project is for you, if you want to write a REST API in Go, and you value the following:

* Design-fist approach using OpenAPI 3
* Type safety: avoid using `context.Context` and `interface{}` as much as possible
* Compiler-driven development: let the generated code guide you what you need to implement
* Approaches like [How I write Go HTTP services after seven years](https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831) but don't want to deal with routing, parameter validation and documentation

## Usage

1. Design your API, e.g. in the [Swagger Editor](https://swagger.io/tools/swagger-editor/)
2. Export the specification as JSON
3. Generate Go code via `go-oas-server generate --file api.json --output-dir=gen`
4. Implement your server and its endpoints:

```
package main

type myServer struct {}

func (s *myServer) Middlewares() gen.Middlewares {
	return Middlewares{}
}

func (s *server) HandlePetIndex() (gen.Middlewares, gen.APIOperationPetIndex) {
	return nil, func(w http.ResponseWriter, r *http.Request, p gen.ParametersPetIndex) {
		// Implement your logic here
	}
}

func main() {
	s := NewOpenAPIServer(&myServer{})
	s.Boot()
	log.Fatal(http.ListenAndServe(":8000", s)
}
```

## Current State

This project is currently only little more than a proof of concept. While the general building blocks are in place, a lot of the API may change. Further, some areas are not even covered yet such as generating support for responses.

Here's a (non-exhaustive) list of what's left to do:
* Support more types and validations
* Endpoint-level Go dependencies
* Responses
* Handling more components
* Make it easier for users to figure out what to implement after code generation
* More testing support
* Better naming support (avoid bad chars, and support [common initialisms](https://github.com/golang/lint/blob/8f45f776aaf18cebc8d65861cc70c33c60471952/lint.go#L771)).

Also, go-oas-server cannot generate code for every possible specification. In part, this limitation exists to avoid complexity.

go-oas-server (currently) does not support:

* Multiple paths with different path parameters in the same position (`/foo/{carId}` and `/foo/{personId}`)
* Parameter serialization.
* Parameters defined with `content` instead of `schema`.
* Links.
* Callbacks.
* Responses (however support for this is planned).

There are further egde cases for sure, if you think you've run into one, please open an issue.


### Why generate code from the specification, and not the other way around?

* Typically, it is faster to define the specification of an API than to implement it.
* If you have access to future consumers of the API endpoint(s) being designed, it is easier to gather feedback based on the specification. Having only code is not a good base for discussion.
* It is non-trivial to control every aspect of the specification from code. Often this involves using lots of annotations in other languages. For example, parameter constraints are easy to generate from specification, but hard to generate from code.
* When the specification is generated, it is often not looked at by the API developer, and therefore not as detailed as it could be.
* Writing API endpoints in Go involves a lot of boilerplate, which is boring to write. A generator fits nicely.

## What is the difference to go-swagger?

* go-swagger is an implementation of Swagger 2.0, not OpenAPI 3.
* go-swagger looks huge and I wanted something smaller that feels more like writing `net/http` handlers.
* oas-go-server does not generate clients.
