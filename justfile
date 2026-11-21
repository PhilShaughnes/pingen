# https://cheatography.com/linux-china/cheat-sheets/justfile/

set dotenv-load := true

# list available recipes
default:
	@just --list

# build the pingen binary
build:
	go build -o pingen .

# run pingen with default settings (source: ./, output: ./public)
run:
	go run .

# run pingen with custom source and output directories
run-custom SOURCE OUTPUT:
	go run . -s {{SOURCE}} -o {{OUTPUT}}

# run unit tests
test:
	go test -v

# build and run the test content
test-content: build
	./pingen -s local/test-content -o local/test-output

# watch and run a go file
watch PATH:
	ls {{PATH}}/* | entr -c go run {{PATH}}/*.go

# watch and run tests for a go file
wtest PATH:
	ls {{PATH}}/* | entr -c go test {{PATH}}/*.go

# build for linux x86_64
build-linux:
	GOOS=linux GOARCH=amd64 go build -o pingen-linux .

# deploy to wyse 3040
deploy: build-linux
	ssh dwys1.lan 'mkdir -p ~/.local/bin'
	rsync -avz pingen-linux wx.lan:~/.local/bin/pingen

# clean build artifacts and generated output
clean:
	rm -f pingen pingen-linux
	rm -rf public/
	rm -rf local/test-output/ dist/

# local release build into dist/ (nothing is published)
snapshot:
	mise exec goreleaser@latest -- goreleaser release --snapshot --clean

# tag pushed main as vVERSION and push the tag; the release workflow publishes it
release VERSION:
	#!/usr/bin/env bash
	set -euo pipefail
	v="{{VERSION}}"; tag="v${v#v}"   # accepts 0.1.0 or v0.1.0
	git fetch -q origin
	main=$(git rev-parse main)
	[ "$main" = "$(git rev-parse origin/main)" ] || { echo "main and origin/main differ - push main first"; exit 1; }
	git tag -a "$tag" -m "$tag" "$main"
	git push origin "$tag"
	echo "tagged $tag - watch: gh run list --workflow release.yml"
