
APP_NAME="Quorum"
APP_VERSION="0.0.2"

.PHONY: start, end, setup clean build-linux build-windows build-mac run test doc coverage sca list

# Package executables for all OS after built-in
all: start setup build-windows build-linux build-mac end

start:
	@echo "Starting build"

end:
	@echo "Build successful"	

setup:
	@rm -rf build/*
	@mkdir -p build/

# Build executable for Windows
build-windows: setup
	@go install github.com/akavel/rsrc@latest
	@rsrc -ico logo.ico -o rsrc.syso
	@GOOS=windows GOARCH=amd64 go build -o build/${APP_NAME}.exe

# Build executable for Linux
build-linux: setup
	@GOOS=linux GOARCH=amd64 go build -o build/linux-${APP_NAME}
	@chmod +x build/linux-${APP_NAME}
	
# Build executable for MacOS
build-mac: setup
	@GOOS=darwin GOARCH=amd64 go build -o build/mac-${APP_NAME}
	@chmod +x build/mac-${APP_NAME}

# Run app
run:
	go run .

# Run test all
test:
	go test ./...

# See doc
doc:
	go doc

# get test coverage
coverage:
	go mod download golang.org/x/tools
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out

# run sca analysis
sca:
	govulncheck ./...


# list all target in makefile
list:
	@grep '^[^#[:space:]].*:' Makefile | grep -v '\.PHONY'