# runs the webapp in development mode, launching client and server separately
go build
sh -c './training-notebook --config='${1}' | tee server.log > /dev/null & npm start --prefix frontend/ | tee frontend.log > /dev/null & wait'
