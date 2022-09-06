#!/bin/bash

USAGE=$(cat <<-END
Usage: 

  ./read_user.sh <user-id>

END
)

curl -H 'Content-Type: application/json' http://localhost:5000/api/v1/users/$1
