#!/bin/bash

echo "package commands" > commands/templates.go
echo "" >> commands/templates.go
echo "const (" >> commands/templates.go

echo -n "templateComponents = \`" >> commands/templates.go
cat templates/components.gogo | sed "s/\`/\` + \"\`\" + \`/g" >> commands/templates.go
echo "\`" >> commands/templates.go

echo -n "templateEndpoints = \`" >> commands/templates.go
cat templates/endpoints.gogo | sed "s/\`/\` + \"\`\" + \`/g" >> commands/templates.go
echo "\`" >> commands/templates.go

echo -n "templateRouting = \`" >> commands/templates.go
cat templates/routing.gogo | sed "s/\`/\` + \"\`\" + \`/g" >> commands/templates.go
echo "\`" >> commands/templates.go

echo -n "templateServer = \`" >> commands/templates.go
cat templates/server.gogo | sed "s/\`/\` + \"\`\" + \`/g" >> commands/templates.go
echo "\`" >> commands/templates.go

echo -n "templateType = \`" >> commands/templates.go
cat templates/type.gogo | sed "s/\`/\` + \"\`\" + \`/g" >> commands/templates.go
echo "\`" >> commands/templates.go

echo -n "templateParameterExtraction = \`" >> commands/templates.go
cat templates/parameter_extraction.gogo | sed "s/\`/\` + \"\`\" + \`/g" >> commands/templates.go
echo "\`" >> commands/templates.go

echo ")" >> commands/templates.go
gofmt -w commands/templates.go
