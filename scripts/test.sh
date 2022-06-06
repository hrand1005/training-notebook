#!/bin/bash


readonly TRAINING_NOTEBOOK="github.com/hrand1005/training-notebook/"
readonly USAGE="Usage: $(basename $0) [-a] [-p package]"

# if [ $#==1 ]; then
#   echo $USAGE
#   exit 1
# fi

while getopts ':ap:h' opt; do
  case "$opt" in
    a)
      echo "Running all tests..."
      go test $TRAINING_NOTEBOOK...
      exit 0
      ;;
    p) 
      echo "Running tests in package '${OPTARG}'"
      go test -timeout 30s $TRAINING_NOTEBOOK/$OPTARG
      exit 0
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

