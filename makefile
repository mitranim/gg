MAKEFLAGS := --silent --always-make
MAKE_CONC := $(MAKE) -j 128 CONF=true clear=$(or $(clear),false)
GO_FLAGS := -tags=$(tags) -mod=mod
VERB := $(if $(filter true,$(verb)),-v,)
FAIL := $(if $(filter false,$(fail)),,-failfast)
SHORT := $(if $(filter true,$(short)),-short,)
CLEAR := $(if $(filter false,$(clear)),,-c)
PROF := $(if $(filter true,$(prof)), -cpuprofile=cpu.prof -memprofile=mem.prof,)
TEST_FLAGS := $(GO_FLAGS) -count=1 $(VERB) $(FAIL) $(SHORT) $(PROF)
TEST := test $(TEST_FLAGS) -timeout=1s -run=$(run)
PKG := ./$(or $(pkg),...)
BENCH := test $(TEST_FLAGS) -run=- -bench=$(or $(run),.) -benchmem -benchtime=256ms
GOW_HOTKEYS := -r=$(if $(filter 0,$(MAKELEVEL)),true,false)
GOW := gow $(CLEAR) $(GOW_HOTKEYS) $(VERB) -e=go,mod,pgsql
WATCH := watchexec -r $(CLEAR) -d=0 -n
DOC_HOST := localhost:58214
OK = echo [$@] ok
TAG := $(or $(and $(ver),v0.1.$(ver)),$(tag))

default: test_w

watch:
	$(MAKE_CONC) test_w lint_w

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
	$(MAKE_CONC) prof_cpu prof_mem

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
	$(MAKE_CONC) test lint

# Examples:
# `make pub ver=1`.
# `make pub tag=v0.0.1`.
pub: prep
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
