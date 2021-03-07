.PHONY:
run:
	redis-server --port 6379 & go run main.go --proto=http


pull:
	sudo docker pull kostikan/proxy:latest

upload:
	sudo docker build -t kostikan/proxy:latest -f ./Dockerfile .
	sudo docker push kostikan/proxy:latest
	sudo APP_VERSION=latest docker-compose up

start:
	sudo APP_VERSION=latest docker-compose up