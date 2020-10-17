SHELL=/bin/bash -o pipefail

#  EXECUTABLES = docker-compose docker node npm go
#  K := $(foreach exec,$(EXECUTABLES),\
#          $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

export GO111MODULE := on
export PATH := .bin:${PATH}

GO_DEPENDENCIES = github.com/ory/go-acc \
				  github.com/sqs/goreturns \
				  github.com/ory/x/tools/listx \
				  github.com/markbates/pkger/cmd/pkger \
				  github.com/golang/mock/mockgen \
				  github.com/go-swagger/go-swagger/cmd/swagger \
				  golang.org/x/tools/cmd/goimports \
				  github.com/mikefarah/yq

define make-go-dependency
  # go install is responsible for not re-building when the code hasn't changed
  .bin/$(notdir $1): go.mod go.sum Makefile
		GOBIN=$(PWD)/.bin/ go install $1
endef
$(foreach dep, $(GO_DEPENDENCIES), $(eval $(call make-go-dependency, $(dep))))
$(call make-lint-dependency)

.bin/clidoc:
		go build -o .bin/clidoc ./cmd/clidoc/.

docs/cli: .bin/clidoc
		clidoc .

.bin/traefik:
		https://github.com/containous/traefik/releases/download/v2.3.0-rc4/traefik_v2.3.0-rc4_linux_amd64.tar.gz \
			tar -zxvf traefik_${traefik_version}_linux_${arch}.tar.gz

.bin/cli: go.mod go.sum Makefile
		go build -o .bin/cli -tags sqlite github.com/ory/cli

node_modules: package.json Makefile
		npm ci

docs/node_modules: docs/package.json
		cd docs; npm ci

.bin/golangci-lint: Makefile
		curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b .bin v1.28.3

.bin/hydra: Makefile
		bash <(curl https://raw.githubusercontent.com/ory/hydra/master/install.sh) -b .bin v1.7.4

.PHONY: docs
docs: docs/node_modules
		cd docs; npm run build

.PHONY: lint
lint: .bin/golangci-lint
		golangci-lint run -v ./...

.PHONY: cover
cover:
		go test ./... -coverprofile=cover.out
		go tool cover -func=cover.out

.PHONY: mocks
mocks: .bin/mockgen
		mockgen -mock_names Manager=MockLoginExecutorDependencies -package internal -destination internal/hook_login_executor_dependencies.go github.com/zzpu/openuser/selfservice loginExecutorDependencies

.PHONY: install
install: .bin/packr2 .bin/pack
		make pack
		GO111MODULE=on go install -tags sqlite .
		packr2 clean

.PHONY: test-resetdb
test-resetdb:
		script/testenv.sh

.PHONY: test
test:
		go test -p 1 -tags sqlite -count=1 -failfast ./...

# Generates the SDK
.PHONY: sdk
sdk: .bin/swagger .bin/cli
		swagger generate spec -m -o .schema/api.swagger.json -x internal/httpclient
		cli dev swagger sanitize ./.schema/api.swagger.json
		swagger validate ./.schema/api.swagger.json
		swagger flatten --with-flatten=remove-unused -o ./.schema/api.swagger.json ./.schema/api.swagger.json
		swagger validate ./.schema/api.swagger.json
		rm -rf internal/httpclient
		mkdir -p internal/httpclient
		swagger generate client -f ./.schema/api.swagger.json -t internal/httpclient -A Ory_Kratos
		make format

.PHONY: quickstart
quickstart:
		docker pull oryd/kratos:latest-sqlite
		docker pull oryd/kratos-selfservice-ui-node:latest
		docker-compose -f quickstart.yml -f quickstart-standalone.yml up --build --force-recreate

.PHONY: quickstart-dev
quickstart-dev:
		docker build -f .docker/Dockerfile-build -t oryd/kratos:latest-sqlite .
		docker-compose -f quickstart.yml -f quickstart-standalone.yml up --build --force-recreate

# Formats the code
.PHONY: format
format: .bin/goreturns
		goreturns -w -local github.com/ory $$(listx .)
		cd docs; npm run format
		npm run format

# Runs tests in short mode, without database adapters
.PHONY: docker
docker:
		docker build -f .docker/Dockerfile-build -t oryd/kratos:latest-sqlite .

.PHONY: test-e2e
test-e2e: node_modules test-resetdb
		source script/test-envs.sh
		test/e2e/run.sh sqlite
		test/e2e/run.sh postgres
		test/e2e/run.sh cockroach
		test/e2e/run.sh mysql

.PHONY: migrations-sync
migrations-sync: .bin/cli
		cli dev pop migration sync persistence/sql/migrations/templates persistence/sql/migratest/testdata

.PHONY: migrations-render
migrations-render: .bin/cli
		cli dev pop migration render persistence/sql/migrations/templates persistence/sql/migrations/sql

.PHONY: migrations-render-replace
migrations-render-replace: .bin/cli
		cli dev pop migration render -r persistence/sql/migrations/templates persistence/sql/migrations/sql

.PHONY: migratest-refresh
migratest-refresh:
		cd persistence/sql/migratest; go test -tags sqlite,refresh -short .

.PHONY: pack
pack: .bin/pkger
		pkger -exclude node_modules -exclude docs -exclude .git -exclude .github -exclude .bin -exclude test -exclude script -exclude contrib
