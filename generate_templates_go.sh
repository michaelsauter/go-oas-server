#!/bin/bash

echo "package commands" > pkg/commands/templates.go
echo "" >> pkg/commands/templates.go
echo "const (" >> pkg/commands/templates.go

echo -n "templateComponents = \`" >> pkg/commands/templates.go
cat pkg/templates/components.gogo | sed "s/\`/\` + \"\`\" + \`/g" >> pkg/commands/templates.go
echo "\`" >> pkg/commands/templates.go

echo -n "templateEndpoints = \`" >> pkg/commands/templates.go
cat pkg/templates/endpoints.gogo | sed "s/\`/\` + \"\`\" + \`/g" >> pkg/commands/templates.go
echo "\`" >> pkg/commands/templates.go

echo -n "templateRouting = \`" >> pkg/commands/templates.go
cat pkg/templates/routing.gogo | sed "s/\`/\` + \"\`\" + \`/g" >> pkg/commands/templates.go
echo "\`" >> pkg/commands/templates.go

echo -n "templateServer = \`" >> pkg/commands/templates.go
cat pkg/templates/server.gogo | sed "s/\`/\` + \"\`\" + \`/g" >> pkg/commands/templates.go
echo "\`" >> pkg/commands/templates.go

echo -n "templateType = \`" >> pkg/commands/templates.go
cat pkg/templates/type.gogo | sed "s/\`/\` + \"\`\" + \`/g" >> pkg/commands/templates.go
echo "\`" >> pkg/commands/templates.go

echo -n "templateParameterExtraction = \`" >> pkg/commands/templates.go
cat pkg/templates/parameter_extraction.gogo | sed "s/\`/\` + \"\`\" + \`/g" >> pkg/commands/templates.go
echo "\`" >> pkg/commands/templates.go

echo ")" >> pkg/commands/templates.go
gofmt -w pkg/commands/templates.go
