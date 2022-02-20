# runs the webapp in development mode, launching client and server separately
sh -c 'go run main.go | tee server.log > /dev/null & npm start --prefix frontend/ | tee frontend.log > /dev/null & wait'
