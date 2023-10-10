.PHONY: docker
docker:
	@rm mybook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o mybook .
	@docker rmi -f lee/mybook-live:v0.0.1
	@docker build -t lee/mybook-live:v0.0.1 .