# make sure that the wasm_exec.js used matches the go version
cp $(go env GOROOT)/misc/wasm/wasm_exec.js frontend/wasm_exec.js
# compile frontend (namely compile go code to wasm)
(cd frontend;GOOS=js GOARCH=wasm go build -o lib.wasm main.go)
# run frontend

