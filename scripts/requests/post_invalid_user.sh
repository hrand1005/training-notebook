#!/bin/bash

USAGE=$(cat <<-END
This script contains multiple invalid post requests.
Edit the script and un-comment the request you want to curl.
If you add new invalid requests to this file, please leave a 
comment explaining the invalid case that they exercise.

END
)

echo "$USAGE"
echo ""

# invalid first name too short
# curl -d '{"data":{"type": "user","attributes":{"first-name":"h","last-name":"rand","email":"herb@yahoo.mail"}}}' -H 'Content-Type: application/json' -X POST http://localhost:5000/api/v1/users
#
# # invalid first name too long
# curl -d '{"data":{"type": "user","attributes":{"first-name":"hehehehehehehehehehehehehehehehehe","last-name":"rand","email":"herb@yahoo.mail"}}}' -H 'Content-Type: application/json' -X POST http://localhost:5000/api/v1/users
#
# # invalid last name too short
# curl -d '{"data":{"type": "user","attributes":{"first-name":"her","last-name":"r","email":"herb@yahoo.mail"}}}' -H 'Content-Type: application/json' -X POST http://localhost:5000/api/v1/users
#
# # invalid last name too long
# curl -d '{"data":{"type": "user","attributes":{"first-name":"her","last-name":"rrarararararararararararararararararararararararararararararaa","email":"herb@yahoo.mail"}}}' -H 'Content-Type: application/json' -X POST http://localhost:5000/api/v1/users
#
# # invalid email too long
# curl -d '{"data":{"type": "user","attributes":{"first-name":"her","last-name":"ran","email":"hherbherbherbherbherbherbherbherbherbherberb@yahoo.mail"}}}' -H 'Content-Type: application/json' -X POST http://localhost:5000/api/v1/users
#
# # invalid email (not an email address)
# curl -d '{"data":{"type": "user","attributes":{"first-name":"her","last-name":"ran","email":"herbyahoo.mail"}}}' -H 'Content-Type: application/json' -X POST http://localhost:5000/api/v1/users
