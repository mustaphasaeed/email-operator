# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
.PHONY: deploy
deploy:
	kubectl apply -f config/default
	kubectl apply -f config/rbac
	kubectl apply -f config/crd
	sed 's|BASE_IMAGE|${IMG}|g' config/manager/manager.yaml | kubectl apply -f -
	sed 's|BASE64_ENCODED_PLACEHOLDER|${MAIL_SEND_TOKEN}|g' config/secrets/mailersend.yaml | kubectl apply -f -

# Build the docker image
.PHONY: docker-build
docker-build:
	docker build -t ${IMG} .

# Push the docker image
.PHONY: docker-push
docker-push:
	docker push ${IMG}

# Clean the build artifacts
.PHONY: clean
clean:
	rm -rf bin

# Define the default target
.PHONY: all
all: build
