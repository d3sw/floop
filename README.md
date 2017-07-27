# floop
floop is a tool to augment any process to provide event based lifecycle calls.  It takes
a process and wraps it to perform user defined actions at different phases of a long
running process

### Development
This project requires go1.8.1+.  It can be built as follows:

```
# Get dependencies
make deps

# Build
make floop
```

#### Example
Example configs can be found under the test-data directory.
