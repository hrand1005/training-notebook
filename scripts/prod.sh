# starts training-notebook webapp in prod mode
# perhaps build things first... go build, yarn build --prefix frontend/
go build
yarn --cwd frontend/ build 
./training-notebook --config=$1 --prod=true