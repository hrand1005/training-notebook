#!/bin/bash

# shortcut for seeding the test database defined by the config in configs/test_config.yaml
go run tools/seeder.go --config=configs/test_config.yaml
