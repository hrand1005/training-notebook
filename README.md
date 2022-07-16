### Requirements

From [gRPC documentation](https://grpc.io/docs/languages/go/quickstart/)
Install protocol compiler and plugins for Go (for development), e.g. on Ubuntu
```
sudo apt install -y protobuf-compiler
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```
Update your PATH so that the protoc compiler can find the plugins:
```
export PATH="$PATH:$(go env GOPATH)/bin"
```