.PHONY: build certs tests

# ------------------------------------------------------------------------
# variables
# ------------------------------------------------------------------------

brewfile:=brewfile
certs_dir:=./certs
dc_name:="dc1"
python_version:="3.9.1"
stack_file:=./ops/local/stack.yaml
debug_config:=./dev/config/consul/debug.yaml
release_config:=./dev/config/consul/release.yaml
current:=$(shell cat .python-version)
repo:=$(shell basename $(CURDIR))

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

--debug:
	@echo
	@docker compose -f "$(stack_file)" up --no-color --quiet-pull --force-recreate -d
	@echo
	@consul kv put "services/customers" @$(debug_config)

--release:
	@echo
	@pack build delineateio/customers -q --builder gcr.io/buildpacks/builder:v1 -p ./dev 1> /dev/null
	@echo
	@docker compose --profile release -f "$(stack_file)" up --scale customers=3 --quiet-pull --force-recreate -d
	@echo
	@consul kv put "services/customers" @$(release_config)

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

debug: --debug ps

release: --release ps

ps:
	@echo
	@docker ps -a --format="table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.State}}\t{{.Networks}}" -a

tests: --clean
	@echo
	@behave tests

down:
	@echo
	@docker compose -f "$(stack_file)" down --remove-orphans

rename: --venv --references
