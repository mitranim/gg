MAKEFLAGS  := --silent --always-make
PAR        := $(MAKE) -j 128
GO_FLAGS   := -tags=$(tags) -mod=mod
VERB       := $(if $(filter $(verb),true), -v,)
SHORT      := $(if $(filter $(short),true), -short,)
PROF       := $(if $(filter $(prof),true), -cpuprofile=cpu.prof -memprofile=mem.prof,)
TEST_FLAGS := $(GO_FLAGS) -count=1 $(VERB) $(SHORT) $(PROF)
TEST       := test $(TEST_FLAGS) -timeout=8s -run=$(run)
FEAT       := ./$(or $(feat),...)
BENCH      := test $(TEST_FLAGS) -run=- -bench=$(or $(run),.) -benchmem -benchtime=128ms
WATCH      := watchexec -r -c -d=0 -n

default: test_w

watch:
	$(PAR) test_w lint_w

test_w:
	gow -c -v $(TEST) $(FEAT)

test:
	go $(TEST) $(FEAT)

bench_w:
	gow -c -v $(BENCH) $(FEAT)

bench:
	go $(BENCH) $(FEAT)

lint_w:
	$(WATCH) -- $(MAKE) lint

lint:
	golangci-lint run
	echo [lint] ok

prof:
	$(PAR) prof_cpu prof_mem

prof_cpu:
	go tool pprof -web cpu.prof

prof_mem:
	go tool pprof -web mem.prof
