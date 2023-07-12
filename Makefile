.PHONY: build_client build_server run_server run_client

DOCKER_CLIENT_IMAGE = wisdom-client
DOCKER_SERVER_IMAGE = wisdom-server


build-server:
	docker build -t $(DOCKER_SERVER_IMAGE) -f Dockerfile.server .

build-client:
	docker build -t $(DOCKER_CLIENT_IMAGE) -f Dockerfile.client .

run-server:
	docker run --rm -p 8080:8080 --name server -v $(PWD)/.env:/app/.env --env-file .env $(DOCKER_SERVER_IMAGE)

run-client:
	docker run --rm --name client --link server $(DOCKER_CLIENT_IMAGE)
