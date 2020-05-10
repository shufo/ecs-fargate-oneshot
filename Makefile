MAKEFLAGS += --always-make
MAKEFLAGS += --silent
MAKEFLAGS += --ignore-errors
MAKEFLAGS += --no-print-directory

# constants
export PROJECT = $(shell basename `pwd`)
export UID = $(shell id -u)
export GID = $(shell id -g)


.PHONY: logs
ifeq (logs,$(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif
logs: ## Display container's log : ## make logs, make logs app
	@docker logs -f $(PROJECT)-$(RUN_ARGS)

.PHONY: up
ifeq (up,$(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif
up: ## Create and start containers : ## make up, make up mysql
	docker-compose -f docker-compose.yml -p $(PROJECT) up -d $(RUN_ARGS)

.PHONY: kill
ifeq (kill,$(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif
kill: ## kill containers : ## make kill, make kill mysql
	docker-compose -f docker-compose.yml -p $(PROJECT) kill $(RUN_ARGS)

.PHONY: test
ifeq (test,$(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif
test: ## test containers : ## make test, make test mysql
	docker-compose -f docker-compose.yml -p $(PROJECT) exec app go test -v ./...

.PHONY: rm
ifeq (rm,$(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif
rm: ## Stop & Remove containers : ## make rm, make rm mysql
	docker-compose -f docker-compose.yml -p $(PROJECT) kill $(RUN_ARGS) && \
	docker-compose -f docker-compose.yml -p $(PROJECT) rm -f $(RUN_ARGS)

ps: ## List containers : ## make ps
	docker-compose -f docker-compose.yml -p $(PROJECT) ps

.PHONY: restart
ifeq (restart,$(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif
restart: ## Restart services : ## make restart, make restart app
	docker-compose -f docker-compose.yml -p $(PROJECT) kill $(RUN_ARGS) && \
	docker-compose -f docker-compose.yml -p $(PROJECT) rm -f $(RUN_ARGS) && \
	docker-compose -f docker-compose.yml -p $(PROJECT) up -d $(RUN_ARGS)

.PHONY: attach
ifeq (attach,$(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif
attach: ## Attach to container : ## make attach app
	docker-compose -f docker-compose.yml -p $(PROJECT) exec -u $(UID):$(UID) $(RUN_ARGS) sh -c "[ -f /bin/bash ] && /bin/bash || /bin/sh"

.PHONY: help
help: ## Show this help message : ## make help
	@echo -e "\nUsage: make [command] [args]\n"
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ": ## "}; {printf "\t\033[36m%-20s\033[0m \033[33m%-30s\033[0m (e.g. \033[32m%s\033[0m)\n", $$1, $$2, $$3}'
	@echo -e "\n"

.DEFAULT_GOAL := help
