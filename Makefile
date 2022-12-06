.PHONY: run build clean deploy-dev

.PHONY: help
help:
	@echo 'Usage:'
	@echo '    make run						Run the program.'
	@echo '    make build					Build the programs.'
	@echo '    make deploy-dev				Deploy the programs into dev environment.'
	@echo '    make clean					Clean the built bin.'
	@echo

.PHONY: build
build:
	GOARCH=arm64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -tags lambda.norpc -o bin/demo/bootstrap services/demo/main.go
	ditto -c -k --sequesterRsrc bin/demo/bootstrap bin/demo/bootstrap.zip

.PHONY: deploy-dev
deploy-dev:
	serverless deploy --stage dev --aws-profile mfa

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: run
run:
	./scripts/run.sh