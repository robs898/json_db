server.key:
	openssl req -x509 -nodes -newkey rsa:2048 -keyout server.key -out server.crt -days 3650 -subj '/CN=robbie.casa'

run: server.key
	go fmt
	sudo docker build -t json_db .
	sudo docker run --rm -p 8443:8443 json_db

clean:
	-rm -rf server.key server.crt
	-sudo docker rmi json_db
