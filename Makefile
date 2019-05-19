build: server.go
	go fmt
	go build server.go

run: server.go
	go fmt
	go run server.go

certs:
	-rm -rf server.key server.crt
	openssl req -x509 -nodes -newkey rsa:2048 -keyout server.key -out server.crt -days 3650

clean:
	-rm -rf json_db
	-rm -rf server.key server.crt

clean-db:
	-rm -rf *.json

docker-build: Dockerfile
	docker build -t json_db .
