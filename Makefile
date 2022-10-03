# Copyright (c) 2021 Acronis International GmbH
#
# Permission is hereby granted, free of charge, to any person obtaining a copy of
# this software and associated documentation files (the "Software"), to deal in
# the Software without restriction, including without limitation the rights to
# use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
# the Software, and to permit persons to whom the Software is furnished to do so,
# subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
# FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
# COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
# IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
# CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
#
# Makefile for sample connector repository
#
# Targets:
# --------
# - lint	Perform linting (requires docker)
# - test	Perform unit tests (requires docker and docker-compose)
# - e2etest Perform end-to-end tests (requires docker and docker-compose)
# - deploy  Run Sample Connector and Sample External System in dockers
#
# Configurable parameters:
# ------------------------

PROJECT_PATH := github.com/acronis
PROJECT_NAME := acronis-cyber-cloud-go-sample-connector
RANDOM_IDENTIFIER := $(shell /bin/bash -c "echo $$RANDOM")

# END of configuration. Change below this line only if modifying logic.
# ------------------------------------------------------------------------------

.env:
	printf '%s\n%s\n%s\n%s\n%s\n%s\n%s' \
	'AUTH_CLIENT_ID=' \
	'AUTH_CLIENT_SECRET=' \
	'DB_USER=postgres' \
	"DB_PASSWORD=$$(openssl rand -hex 16)" \
	'SSO_AUTH_CLIENT_ID=' \
	'SSO_AUTH_CLIENT_SECRET=' \
	'SSO_AUTH_SESSION_SECRET=' \
	> .env

lint:
	@echo "--------------------------------------------------------------------------------"
	@echo "Run linters for Go code"
	@echo "================================================================================"
	docker run --rm -v `pwd`:/go/src/$(PROJECT_PATH)/$(PROJECT_NAME) -w /go/src/$(PROJECT_PATH)/$(PROJECT_NAME) golangci/golangci-lint:v1.31.0 golangci-lint run -v --timeout 5m0s

test: .env
test: export COMPOSE_FLAG := -f docker-compose.yaml -f docker-compose.ci.yaml
test:
	@echo "--------------------------------------------------------------------------------"
	@echo "Run tests"
	@echo "================================================================================"
	
	docker-compose -p $(RANDOM_IDENTIFIER) $(COMPOSE_FLAG) down --remove-orphans
	docker-compose $(COMPOSE_FLAG) -p $(RANDOM_IDENTIFIER) up --scale externalsystem=0 -d --build
	echo 'docker compose up successfully'

	sleep 20s
	docker exec $(RANDOM_IDENTIFIER)_sampleconnector_1 \
		bash -c 'go test -p 1 ./...'\
		|| (echo "Autotest Failed. Printing docker logs..."; \
		docker-compose -p $(RANDOM_IDENTIFIER) logs; \
		docker-compose $(COMPOSE_FLAG) -p $(RANDOM_IDENTIFIER) down --remove-orphans; \
		exit 1) 
	docker-compose $(COMPOSE_FLAG) -p $(RANDOM_IDENTIFIER) down --remove-orphans

e2etest:
	@echo "--------------------------------------------------------------------------------"
	@echo "Run e2e tests"
	@echo "================================================================================"
	docker-compose -f docker-compose.yaml up -d --build
	docker-compose -f docker-compose.yaml -f docker-compose.e2etest.yaml build e2etest
	@echo "Running tests..."
	docker-compose -f docker-compose.yaml -f docker-compose.e2etest.yaml run --rm e2etest /bin/bash -c "go test -v --tags=e2e"
	docker-compose -f docker-compose.yaml down

deploy:
	@echo "--------------------------------------------------------------------------------"
	@echo "Run Connector and Sample External System"
	@echo "================================================================================"
	docker-compose up -d --build

stop:
	@echo "--------------------------------------------------------------------------------"
	@echo "Stop Connector and Sample External System"
	@echo "================================================================================"
	docker-compose stop
