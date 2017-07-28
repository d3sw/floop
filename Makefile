
NAME = floop
FILES = cmd/main.go cmd/cli.go
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
COMMIT = $(shell git rev-parse --short HEAD)
VERSION = $(shell git describe | sed -e "s/^v//")
BUILDTIME = $(shell date +%Y-%m-%dT%T%z)
BUILD_CMD = CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo \
	-ldflags="-X main.branch=${BRANCH} -X main.commit=${COMMIT} -X main.buildtime=${BUILDTIME} -w"

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

	$(eval OS := linux) $(eval OUTFILE := ./dist/$(NAME)-$(OS)-$(VERSION))
	GOOS=$(OS) $(BUILD_CMD) -o $(OUTFILE) $(FILES)
	tar -czf $(OUTFILE).tgz $(OUTFILE); rm -f $(OUTFILE)

	$(eval OS := darwin) $(eval OUTFILE := ./dist/$(NAME)-$(OS)-$$(VERSION))
	GOOS=$(OS) $(BUILD_CMD) -o $(OUTFILE) $(FILES)
	tar -czf $(OUTFILE).tgz $(OUTFILE); rm -f $(OUTFILE)

	$(eval OS := windows) $(eval OUTFILE := ./dist/$(NAME)-$(OS)-$$(VERSION))
	GOOS=$(OS) $(BUILD_CMD) -o $(OUTFILE).exe $(FILES)
	zip $(OUTFILE).zip $(OUTFILE).exe; rm -f $(OUTFILE).exe