
NAME = floop
FILES = child.go config.go main.go writer.go


clean:
	go clean -i ./...
	rm -f $(NAME)

build:
	go build -o $(NAME) .
