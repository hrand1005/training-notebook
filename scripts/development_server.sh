#!/bin/bash

go build cmd/training-notebook/main.go
./main --config=configs/test_config.yaml

