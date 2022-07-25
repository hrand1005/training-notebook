# compile frontend (namely compile go code to wasm)
(cd frontend;GOOS=js GOARCH=wasm go build -o lib.wasm main.go)
# run frontend

