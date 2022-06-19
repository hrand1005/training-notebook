#!/bin/bash

# go project root
readonly TRAINING_NOTEBOOK="github.com/hrand1005/training-notebook/"

# usage where a is all unit tests and i is integration tests
readonly USAGE="Usage: $(basename $0) [-a] [-i] [-p package]"

# test server startup utilities
readonly TEST_CONFIG_PATH="configs/test_config.yaml"
readonly START_TEST_SERVER="sh scripts/dev.sh $TEST_CONFIG_PATH"
readonly TEST_SERVER_HOST="localhost"
readonly TEST_SERVER_ADDR=8080
readonly MAX_RETRIES=10

wait_for_test_server() {
  TRIES=0
  while [ $((TRIES++)) -lt $MAX_RETRIES ]; do
    nc -zv $TEST_SERVER_HOST $TEST_SERVER_ADDR
    if [ $? -eq 0 ]; 
    then
      break
    fi
    sleep 2
  done
}

while getopts ':aip:h' opt; do
  case "$opt" in
    a)
      echo "Running all unit tests..."
      go test -v $(go list $TRAINING_NOTEBOOK... | grep -v /test)
      if [ $? -eq 0 ]; then
        exit 0
      fi
      exit $?
      ;;
    p) 
      echo "Running tests in package '${OPTARG}'"
      go test -v -timeout 30s $TRAINING_NOTEBOOK/$OPTARG
      if [ $? -eq 0 ]; then
        exit 0
      fi
      exit $?
      ;;
    i)
      echo "Running integration tests..."
      $START_TEST_SERVER&
      wait_for_test_server
      go test -v $TRAINING_NOTEBOOK/test
      if [ $? -eq 0 ]; then
        pkill -P $$
        exit 0
      fi
      pkill -P $$
      exit $?
      ;;
    h)
      echo $USAGE
      exit 0
      ;;
    :)
      echo -e "Option requires argument.\n$USAGE"
      exit 1
      ;;
    ?) 
      echo -e "Invalid option.\n$USAGE"
      exit 1
      ;;
  esac
done

shift "$((OPTIND -1))"

