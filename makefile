.PHONY: build certs tests

# ------------------------------------------------------------------------
# variables
# ------------------------------------------------------------------------

brewfile:=brewfile
certs_dir:=./certs
dc_name:="dc1"
python_version:="3.9.1"
stack_file:=./ops/local/stack.yaml
services_config:=./dev/config/consul/services.yaml
up_config:=./dev/config/consul/up.yaml
current:=$(shell cat .python-version)
repo:=$(shell basename $(CURDIR))
rebuild:=false
scale:=3

# ------------------------------------------------------------------------
# private targets
# ------------------------------------------------------------------------

--install:
	@echo
	@brew bundle install --file=$(brewfile)

--venv:
	@echo
	@pyenv install $(python_version) -s
	@pyenv virtualenv -q -f $(python_version) $(repo) 1> /dev/null
	@pyenv local $(repo)
	@pip install -q --upgrade pip
	@pip install -qr requirements.txt

--certs:
	@echo
	@rm -rf $(certs_dir)
	@mkdir -p $(certs_dir)
	@cd certs; consul tls ca create
	@cd certs; consul tls cert create -server -dc $(dc_name)

--hooks:
	@echo
	@pre-commit install
	@pre-commit autoupdate
	@echo
	@git add .
	@-pre-commit

--services:
	@echo
	@docker compose -f "$(stack_file)" up \
		--no-color \
		--no-recreate \
		--quiet-pull \
		--detach
	@echo
	@consul kv put "services/customers" @$(services_config)

--up:
    ifeq ($(rebuild), true)
		@echo
		@pack build delineateio/customers -q --builder gcr.io/buildpacks/builder:v1 -p ./dev 1> /dev/null
    endif
	@echo
	@docker compose --profile full -f "$(stack_file)" up \
		--no-color \
		--scale customers=$(scale) \
		--quiet-pull \
		--detach
	@echo
	@consul kv put "services/customers" @$(up_config)

--clean:
	@psql -h "localhost" -U "postgres" -c "DELETE FROM customers;" 1> /dev/null
	@rm -rf build

--references:
	@sed -i '' s/$(current)/$(repo)/g readme.md
	@sed -i '' s/$(current)/$(repo)/g .python-version

# ------------------------------------------------------------------------
# public targets
# ------------------------------------------------------------------------

list:
	@echo
	@cat brewfile
	@echo
	@cat requirements.txt

init: --install --certs --venv --hooks

graph:
	@cd ./dev; go mod tidy
	@cd ./dev; go mod graph | tee ../graph.txt

build:
	@echo
	@cd ./dev/src; go build -o ../../build/customers

services: --services ps

up: --up ps

ps:
	@echo
	@docker ps \
		--format="table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.State}}\t{{.Networks}}" \
		--filter 'network=local_consul' \
		--all

tests: --clean
	@echo
	@behave tests

down:
	@echo
	@docker compose -f "$(stack_file)" down --remove-orphans

rename: --venv --references
