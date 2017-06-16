
NAME = floop
FILES = cmd/main.go cmd/cli.go

BUILD_CMD = go build

clean:
	go clean -i ./...
	rm -f $(NAME)
	rm -rf ./dist

deps:
	go get -d ./...

$(NAME):
	go build -o $(NAME) $(FILES)

dist:
	[ -d ./dist ] || mkdir ./dist
	GOOS=linux $(BUILD_CMD) -o ./dist/$(NAME)-linux $(FILES)
	GOOS=darwin $(BUILD_CMD) -o ./dist/$(NAME)-darwin $(FILES)
	GOOS=windows $(BUILD_CMD) -o ./dist/$(NAME) $(FILES)
