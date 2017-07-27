# Contributing to floop

The purpose of this repository is to continue to evolve floop, adding additional integration protocols
and configuration options.  All development is done publicly on GitHub, and we look forward to working
with a community of both users and developers to make floop even more powerful and intuitive.

We are following the GitFlow workflow. The active development branch is 
[develop](https://github.com/d3sw/floop/tree/develop), the stable branch is 
[master](https://github.com/d3sw/floop/tree/master).

Contributions will be accepted to the [develop](https://github.com/d3sw/floop/tree/develop) only.

## How to submit an issue for a feature or bug

Your issue may already be reported! Please search existing floop [issues](https://github.com/d3sw/floop/issues) 
before creating one. Be ready to provide the following information.
* Expected Behavior
* Current Behavior
* Possible Solution
* Steps to Reproduce
* Any Additional Context
* Environment Details

## How to submit a proposed code change

1. Please create an [Issue](https://github.com/d3sw/floop/issues) to discuss 
the details of the feature with the community as well as with the repository administrators.

2. Once the issue has been discussed and a concensus has been reached on the implementation, please 
follow the procedures below to submit your change request.

    a. Fork floop on GitHub. [How to Fork a Repository](http://help.github.com/fork-a-repo/)

    b. Create a new branch on your fork.

    c. Push your code change to your new branch

    d. Initiate a pull request on github [How to Send a Pull Request](http://help.github.com/send-pull-requests/)

## Development Requirements

This project requires go1.8.1+.

### Go Get Dependencies
```bash
go get -d ./...
# or
make deps
```

### Build Binaries
```bash
go build -o floop cmd/main.go cmd/cli.go
# or
make floop
```
