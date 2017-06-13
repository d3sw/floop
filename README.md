# floop


### Development

#### Example

```
go run ./cmd/*.go \
    docker run --rm  opencoconut/ffmpeg \
    -progress /dev/stdout \
    -i http://files.coconut.co.s3.amazonaws.com/test.mp4 \
    -f webm -c:v libvpx -c:a libvorbis \
    test.webm
```

### TODO

- [ ] Service discovery
- [ ] Update compile options
