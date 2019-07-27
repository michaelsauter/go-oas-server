#!/usr/bin/env bash
set -eux

version=$1

if [ -z "$version" ]; then
  echo "No version passed! Example usage: ./release.sh 1.0.0"
  exit 1
fi

echo "Running tests..."
make test

echo "Update version..."
old_version=$(grep -o "[0-9]*\.[0-9]*\.[0-9]*" cmd/go-oas-server/main.go)
sed -i.bak 's/fmt.Println("'$old_version'+master")/fmt.Println("'$version'")/' cmd/go-oas-server/main.go
sed -i.bak 's/'$old_version'/'$version'/' README.md

echo "Mark version as released in changelog..."
today=$(date +'%Y-%m-%d')
sed -i.bak 's/Unreleased/Unreleased\
\
## ['$version'] - '$today'/' CHANGELOG.md

rm *.bak

echo "Build binaries..."
make build

echo "Update repository..."
git add cmd/go-oas-server/main.go README.md CHANGELOG.md
git commit -m "Bump version to ${version}"
git tag --message="v$version" --force "v$version"
git tag --message="latest" --force latest

echo "Set master version again"
sed -i.bak 's/fmt.Println("'$version'")/fmt.Println("'$version'+master")/' cmd/go-oas-server/main.go
rm cmd/go-oas-server/main.go.bak
git add cmd/go-oas-server/main.go
git commit -m "Set master version to ${version}+master"

echo "v$version tagged."
echo "Now, run 'git push origin master && git push --tags --force' and publish the release on GitHub."
