#!/bin/bash

source test/behave/.venv/training-notebook/bin/activate
behave test/behave/features -D server_config=${1} --no-capture
deactivate
