run:
	go run cmd/main.go

docker-build:
	docker build --no-cache -t uartweb/proxy .

docker-run:
	echo "Running on port 8000"
	docker run --rm -d -p 8000:8000 --net uart_net --name uart-proxy uartweb/proxy

docker-stop:
	docker stop uart-proxy

docker-remove:
	make docker-stop
	docker rm uart-proxy

docker-push:
	make docker-build
	docker push uartweb/proxy:latest

docker-deploy:
	make docker-build
	make docker-run

lint:
	golangci-lint run

gomod-download:
	go get -u github.com/gorilla/handlers
	go get -u github.com/gorilla/mux
	go get -u github.com/spf13/viper

gomod-tidy:
	go mod tidy -go=1.19 -compat=1.19

gomod-update:
	make gomod-download
	make gomod-tidy
