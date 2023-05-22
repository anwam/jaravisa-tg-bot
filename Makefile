.PHONY: all build tag push

all: build tag push

build:
	docker build --platform linux/amd64 -t jaravisa-tg-bot:latest .

tag:
	docker tag jaravisa-tg-bot:latest asia-southeast1-docker.pkg.dev/swerbillionaire-atlas/cloud-run-source-deploy/jaravisa-tg-bot:latest

push:
	docker push asia-southeast1-docker.pkg.dev/swerbillionaire-atlas/cloud-run-source-deploy/jaravisa-tg-bot:latest
