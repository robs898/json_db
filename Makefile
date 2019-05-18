build: server.go
	go fmt
	go build server.go

run: server.go
	go fmt
	go run server.go

clean:
	-rm -rf json_db

clean-db:
	-rm -rf *.json

docker-build: Dockerfile
	docker build -t json_db .
