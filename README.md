# example_api
Sample API server as a container

## to build it
go build -o bin/api_server cmd/api_server/*.go

## to build smaller
go build -ldflags="-s -w" -o bin/api_server cmd/api_server/*.go


## to run it
$ bin/api_server


## API list
/version/ returns json of version, subversion and codename
/libs/ returns json of all the libraries
/lib/<uuid> returns content of library with uuid <uuid>


## Dockerize it
docker build --tag docker-api-server .

# References
https://boyter.org/posts/how-to-start-go-project-2023/
 and discussion https://news.ycombinator.com/item?id=36046662
 
https://docs.docker.com/language/golang/build-images/ 
