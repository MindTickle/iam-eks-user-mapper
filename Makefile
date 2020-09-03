.PHONY: build deploy

DOCKER_IMAGE='191195949309.dkr.ecr.ap-southeast-1.amazonaws.com/devops/iam-eks-user-mapper:latest'

docker_build:
	docker build -t ${DOCKER_IMAGE} .

docker_push:
	aws ecr get-login-password --region ap-southeast-1 | docker login --username AWS --password-stdin 191195949309.dkr.ecr.ap-southeast-1.amazonaws.com
	docker push ${DOCKER_IMAGE}

deploy:
	kubectl apply -f kubernetes/
