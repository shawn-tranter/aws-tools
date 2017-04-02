ENV_VAR="GOPATH=`pwd`"

all: build

random-instance-killer: src/random-instance-killer/*.go
	GOPATH=`pwd` go build random-instance-killer

build: random-instance-killer

get:
	GOPATH=`pwd` go get github.com/aws/aws-sdk-go/aws
