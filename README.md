# literate-bassoon
## Docker
> docker-compose -f docker-compose-single-broker.yml up -d
## Influx account
> Go to localhost:8087 (Chronograf) and create new database and user
## CLI
### Start server
> go run api.go
### Start consumer
> go run reader/reader.go
