docker build -t cons-alpha:red -f Consumer.dockerfile .
docker image tag cons-alpha:red controllerbase/cons-alpha:red
docker build -t prod-alpha:red -f Server.dockerfile .
docker image tag prod-alpha:red controllerbase/prod-alpha:red

docker push controllerbase/cons-alpha:red
docker push controllerbase/prod-alpha:red