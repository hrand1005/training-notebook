# script is intended to be run from project root

# --go_out for pb.go file, go-grpc_out for _grpc.pb.go file
protoc --proto_path=set set/set.proto --go_out=. --go-grpc_out=. --go_opt=paths=source_relative