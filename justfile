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
	rsync -avz pingen-linux dwys1.lan:~/.local/bin/pingen

# clean build artifacts and generated output
clean:
	rm -f pingen pingen-linux
	rm -rf public/
	rm -rf local/test-output/
