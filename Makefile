
NAME = floop
FILES = cmd/main.go cmd/cli.go
COMMIT = $(shell git rev-parse --short HEAD)
VERSION = $(shell grep 'const VERSION' version.go | cut -d'"' -f 2)
BUILDTIME = $(shell date +%Y-%m-%dT%T%z)
BUILD_CMD = CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo \
	-ldflags="-X main.commit=${COMMIT} -X main.buildtime=${BUILDTIME} -w"
GOOS = $(shell go env GOOS)
BUILDLOC = ./dist

version:
	@echo $(VERSION)

clean:
	go clean -i ./...
	rm -rf $(NAME) ./dist coverage.out

deps:
	go get -d ./...

test:
	go test -v -coverprofile=coverage.out .

${NAME}:
	$(BUILD_CMD) -o $(NAME) $(FILES)

dist: clean
	[ -d ./dist ] || mkdir ./dist

	$(BUILD_CMD) -o $(BUILDLOC)/$(NAME) $(FILES)
	tar -czf $(BUILDLOC)/$(NAME)-$(GOOS)-$(VERSION).tgz $(BUILDLOC)/$(NAME); rm -f $(BUILDLOC)/$(NAME)

	GOOS=linux $(BUILD_CMD) -o $(BUILDLOC)/$(NAME) $(FILES)
	tar -czf $(BUILDLOC)/$(NAME)-linux-$(VERSION).tgz $(BUILDLOC)/$(NAME); rm -f $(BUILDLOC)/$(NAME)

	GOOS=windows $(BUILD_CMD) -o $(BUILDLOC)/$(NAME).exe $(FILES)
	zip $(BUILDLOC)/$(NAME)-windows-$(VERSION).zip $(BUILDLOC)/$(NAME).exe; rm -f $(BUILDLOC)/$(NAME).exe
