# runs the webapp in development mode, launching client and server separately
# uses a hard-coded 'test_config.yaml'
./scripts/backend.sh configs/test_config.yaml % frontend.sh

# Command for running front end and backend together in development mode:
# sh -c './training-notebook --config='${1}' | tee server.log > /dev/null & npm start --prefix frontend/ | tee frontend.log > /dev/null & wait'
