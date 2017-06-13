
NAME = floop
FILES = cmd/main.go

BUILD_CMD = go build

clean:
	go clean -i ./...
	rm -f $(NAME)
	rm -rf ./dist

build:
	go build -o $(NAME) ./cmd/*.go

dist:
	[ -d ./dist ] || mkdir ./dist
	GOOS=linux $(BUILD_CMD) -o ./dist/$(NAME)-linux $(FILES)
	GOOS=darwin $(BUILD_CMD) -o ./dist/$(NAME)-darwin $(FILES)
	GOOS=windows $(BUILD_CMD) -o ./dist/$(NAME) $(FILES)
