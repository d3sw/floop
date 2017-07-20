
NAME = floop
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
COMMIT = $(shell git rev-parse --short HEAD)
BUILDTIME = $(shell date +%Y-%m-%dT%T%z)
BUILD_CMD = CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo \
	-ldflags="-X main.branch=${BRANCH} -X main.commit=${COMMIT} -X main.buildtime=${BUILDTIME} -w"

clean:
	go clean -i ./...
	rm -f $(NAME)
	rm -rf ./dist

deps:
	go get -d ./...

$(NAME):
	go build -o $(NAME) .

# APP_VERSION="v0.1.0" make dist
dist: clean
	[ -d ./dist ] || mkdir ./dist

	GOOS=darwin $(BUILD_CMD) -o ./dist/$(NAME) .
	tar -czf ./dist/$(NAME)-darwin-$(APP_VERSION).tgz ./dist/$(NAME); rm -f ./dist/$(NAME)

	GOOS=linux $(BUILD_CMD) -o ./dist/$(NAME) .
	tar -czf ./dist/$(NAME)-linux-$(APP_VERSION).tgz ./dist/$(NAME); rm -f ./dist/$(NAME)

	GOOS=windows $(BUILD_CMD) -o ./dist/$(NAME).exe .
	zip ./dist/$(NAME)-windows-$(APP_VERSION).zip ./dist/$(NAME).exe; rm -f ./dist/$(NAME).exe

publish:
	chmod a+x publish.sh; ./publish.sh	