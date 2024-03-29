MAKEFLAGS := --silent --always-make
MAKE_PAR := $(MAKE) -j 128
GO_FLAGS := -tags=$(tags) -mod=mod
VERB := $(if $(filter $(verb),true),-v,)
FAIL := $(if $(filter $(fail),false),,-failfast)
SHORT := $(if $(filter $(short),true),-short,)
CLEAR := $(if $(filter $(clear),false),,-c)
PROF := $(if $(filter $(prof),true), -cpuprofile=cpu.prof -memprofile=mem.prof,)
TEST_FLAGS := $(GO_FLAGS) -count=1 $(VERB) $(FAIL) $(SHORT) $(PROF)
TEST := test $(TEST_FLAGS) -timeout=1s -run=$(run)
PKG := ./$(or $(pkg),...)
BENCH := test $(TEST_FLAGS) -run=- -bench=$(or $(run),.) -benchmem -benchtime=256ms
GOW := gow $(CLEAR) -v -e=go,mod,pgsql
WATCH := watchexec -r $(CLEAR) -d=0 -n
DOC_HOST := localhost:58214
OK = echo [$@] ok
TAG := $(or $(and $(ver),v0.1.$(ver)),$(tag))

default: test_w

watch:
	$(MAKE_PAR) test_w lint_w

test_w:
	$(GOW) $(TEST) $(PKG)

test:
	go $(TEST) $(PKG)

bench_w:
	$(GOW) $(BENCH) $(PKG)

bench:
	go $(BENCH) $(PKG)

lint_w:
	$(WATCH) -- $(MAKE) lint

lint:
	golangci-lint run
	$(OK)

vet_w:
	$(WATCH) -- $(MAKE) vet

vet:
	go vet $(GO_FLAGS) $(PKG)
	$(OK)

prof:
	$(MAKE_PAR) prof_cpu prof_mem

prof_cpu:
	go tool pprof -web cpu.prof

prof_mem:
	go tool pprof -web mem.prof

# Requires `pkgsite`:
#   go install golang.org/x/pkgsite/cmd/pkgsite@latest
doc:
	$(or $(shell which open),echo) http://$(DOC_HOST)/github.com/mitranim/gg
	pkgsite $(if $(GOREPO),-gorepo=$(GOREPO)) -http=$(DOC_HOST)

prep:
	$(MAKE_PAR) test lint

# Examples:
# `make release ver=1`.
# `make release tag=v0.0.1`.
release: prep
ifeq ($(TAG),)
	$(error missing tag)
endif
	git pull --ff-only
	git show-ref --tags --quiet "$(TAG)" || git tag "$(TAG)"
	git push origin $$(git symbolic-ref --short HEAD) "$(TAG)"

# Assumes MacOS and Homebrew.
deps:
	go install github.com/mitranim/gow@latest
	brew install -q watchexec golangci-lint
