.PHONY: build deploy

DOCKER_IMAGE='191195949309.dkr.ecr.ap-southeast-1.amazonaws.com/devops/iam-eks-user-mapper:latest'

docker_build:
	docker build -t ${DOCKER_IMAGE} .

docker_push:
	aws ecr get-login --no-include-email | sh -
	docker push ${DOCKER_IMAGE}

deploy:
	kubectl apply -f kubernetes/
