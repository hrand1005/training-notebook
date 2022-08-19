#!/bin/bash

behave test/behave/features -D server_config=${1} --no-capture
