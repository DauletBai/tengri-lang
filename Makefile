# ---------------------------------------------
# Tengri-Lang — Minimal Makefile (AOT + benches + CSV + plots)
# ---------------------------------------------

GO     := go
CC     := cc
CARGO  := cargo

BIN_DIR      := .bin
RUNTIME_DIR  := internal/aotminic/runtime
CBENCH_DIR   := c_benches
RS_DIR       := rust_benches

CFLAGS := -O3 -march=native
RSFLAGS := --release

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

# ---------- Tools (AOT transpiler only) ----------
BIN_AOT := $(BIN_DIR)/tengri-aot

.PHONY: build aot clean
build: aot sort_go sort_c sort_rs aot-examples

aot: | $(BIN_DIR)
	@echo "[build] aot -> $(BIN_AOT)"
	$(GO) build -o $(BIN_AOT) ./cmd/tengri-aot

clean:
	@rm -rf $(BIN_DIR)
	@echo "[clean] $(BIN_DIR) removed."

# ---------- AOT demos ----------
FIB_ITER_TGR := benchmarks/src/fib_iter/tengri/fib_iter_cli.tgr
FIB_REC_TGR  := benchmarks/src/fib_rec/tengri/fib_rec_cli.tgr
SORT_TGR_QS  := benchmarks/src/sort/tengri/sort_cli.tgr
SORT_TGR_MS  := benchmarks/src/sort/tengri/sort_cli_m.tgr

BIN_AOT_ITER := $(BIN_DIR)/fib_cli
BIN_AOT_REC  := $(BIN_DIR)/fib_rec_cli
BIN_AOT_SQS  := $(BIN_DIR)/sort_cli_qsort
BIN_AOT_SMS  := $(BIN_DIR)/sort_cli_msort

$(BIN_AOT_ITER): $(BIN_AOT) $(RUNTIME_DIR)/runtime.c $(RUNTIME_DIR)/runtime.h $(FIB_ITER_TGR) | $(BIN_DIR)
	@echo "[aot] fib_iter -> $@"
	@$(BIN_AOT) -o $(BIN_DIR)/fib_cli.c $(FIB_ITER_TGR)
	$(CC) $(CFLAGS) -I$(RUNTIME_DIR) $(BIN_DIR)/fib_cli.c $(RUNTIME_DIR)/runtime.c -o $@

$(BIN_AOT_REC): $(BIN_AOT) $(RUNTIME_DIR)/runtime.c $(RUNTIME_DIR)/runtime.h $(FIB_REC_TGR) | $(BIN_DIR)
	@echo "[aot] fib_rec -> $@"
	@$(BIN_AOT) -o $(BIN_DIR)/fib_rec_cli.c $(FIB_REC_TGR)
	$(CC) $(CFLAGS) -I$(RUNTIME_DIR) $(BIN_DIR)/fib_rec_cli.c $(RUNTIME_DIR)/runtime.c -o $@

$(BIN_AOT_SQS): $(BIN_AOT) $(RUNTIME_DIR)/runtime.c $(RUNTIME_DIR)/runtime.h $(SORT_TGR_QS) | $(BIN_DIR)
	@echo "[aot] sort qsort -> $@"
	@$(BIN_AOT) -o $(BIN_DIR)/sort_cli_qs.c $(SORT_TGR_QS)
	$(CC) $(CFLAGS) -I$(RUNTIME_DIR) $(BIN_DIR)/sort_cli_qs.c $(RUNTIME_DIR)/runtime.c -o $@

$(BIN_AOT_SMS): $(BIN_AOT) $(RUNTIME_DIR)/runtime.c $(RUNTIME_DIR)/runtime.h $(SORT_TGR_MS) | $(BIN_DIR)
	@echo "[aot] sort mergesort -> $@"
	@$(BIN_AOT) -o $(BIN_DIR)/sort_cli_ms.c $(SORT_TGR_MS)
	$(CC) $(CFLAGS) -I$(RUNTIME_DIR) $(BIN_DIR)/sort_cli_ms.c $(RUNTIME_DIR)/runtime.c -o $@

.PHONY: aot-examples
aot-examples: $(BIN_AOT_ITER) $(BIN_AOT_REC) $(BIN_AOT_SQS) $(BIN_AOT_SMS)

# ---------- Native baselines: sort ----------
BIN_SORT_GO := $(BIN_DIR)/sort_go
BIN_SORT_C  := $(BIN_DIR)/sort_c
BIN_SORT_RS := $(BIN_DIR)/sort_rs

sort_go: | $(BIN_DIR)
	@echo "[build] sort_go -> $(BIN_SORT_GO)"
	$(GO) build -o $(BIN_SORT_GO) ./benchmarks/src/sort/go

sort_c: | $(BIN_DIR)
	@echo "[build] sort_c -> $(BIN_SORT_C)"
	$(CC) $(CFLAGS) benchmarks/src/sort/c/sort.c -o $(BIN_SORT_C)

# nuke broken Cargo.lock (macOS issue) so cargo regenerates it
sort_rs: | $(BIN_DIR)
	@echo "[build] sort_rs -> $(BIN_SORT_RS)"
	@rm -f benchmarks/src/sort/rust/Cargo.lock
	@cd benchmarks/src/sort/rust && $(CARGO) build $(RSFLAGS)
	@cp benchmarks/src/sort/rust/target/release/sort_rs $(BIN_SORT_RS)

# ---------- Baselines: fibonacci ----------
BIN_FIB_ITER_C := $(BIN_DIR)/fib_iter_c
BIN_FIB_REC_C  := $(BIN_DIR)/fib_rec_c
$(BIN_FIB_ITER_C): $(CBENCH_DIR)/fib_iter.c $(CBENCH_DIR)/runtime_cbench.h | $(BIN_DIR)
	@echo "[c] fib_iter -> $@"
	$(CC) $(CFLAGS) $(CBENCH_DIR)/fib_iter.c -o $@
$(BIN_FIB_REC_C): $(CBENCH_DIR)/fib_rec.c $(CBENCH_DIR)/runtime_cbench.h | $(BIN_DIR)
	@echo "[c] fib_rec -> $@"
	$(CC) $(CFLAGS) $(CBENCH_DIR)/fib_rec.c -o $@

BIN_FIB_ITER_RS := $(BIN_DIR)/fib_iter_rs
BIN_FIB_REC_RS  := $(BIN_DIR)/fib_rec_rs
$(BIN_FIB_ITER_RS): | $(BIN_DIR)
	@echo "[rust] fib_iter -> $@"
	@cd $(RS_DIR)/fib_iter && $(CARGO) build $(RSFLAGS)
	@cp $(RS_DIR)/fib_iter/target/release/fib_iter $@
$(BIN_FIB_REC_RS): | $(BIN_DIR)
	@echo "[rust] fib_rec -> $@"
	@cd $(RS_DIR)/fib_rec && $(CARGO) build $(RSFLAGS)
	@cp $(RS_DIR)/fib_rec/target/release/fib_rec $@

# ---------- Bench runners + CSV ----------
SIZE ?= 100000
BENCH_REPS ?= 3

RESULT_DIR := benchmarks/results
DATE_TAG   := $(shell date +"%Y%m%d_%H%M%S")
CSV_PATH   := $(RESULT_DIR)/suite_$(DATE_TAG).csv

.PHONY: bench_sort_all bench_sort_all_aot bench_fib_all bench_all
bench_sort_all: sort_go sort_c sort_rs
	@/usr/bin/time -p env SIZE=$(SIZE) BENCH_REPS=$(BENCH_REPS) $(BIN_SORT_GO)
	@/usr/bin/time -p env SIZE=$(SIZE) BENCH_REPS=$(BENCH_REPS) $(BIN_SORT_C)
	@/usr/bin/time -p env SIZE=$(SIZE) BENCH_REPS=$(BENCH_REPS) $(BIN_SORT_RS)

bench_sort_all_aot: bench_sort_all aot-examples
	@/usr/bin/time -p env BENCH_REPS=$(BENCH_REPS) $(BIN_AOT_SQS) $(SIZE)
	@/usr/bin/time -p env BENCH_REPS=$(BENCH_REPS) $(BIN_AOT_SMS) $(SIZE)

bench_fib_all: aot-examples $(BIN_FIB_ITER_C) $(BIN_FIB_REC_C) $(BIN_FIB_ITER_RS) $(BIN_FIB_REC_RS)
	@/usr/bin/time -p env BENCH_REPS=$(BENCH_REPS) $(BIN_AOT_ITER) 45
	@/usr/bin/time -p env BENCH_REPS=$(BENCH_REPS) $(BIN_AOT_REC) 35
	@/usr/bin/time -p $(BIN_FIB_ITER_C) 45
	@/usr/bin/time -p $(BIN_FIB_REC_C) 35
	@/usr/bin/time -p $(BIN_FIB_ITER_RS) 45
	@/usr/bin/time -p $(BIN_FIB_REC_RS) 35

bench_all: bench_sort_all_aot bench_fib_all | $(RESULT_DIR)
	@echo "impl,task,n,reps,time_ns_avg,checksum_or_result" > $(CSV_PATH)
	@{ \
		(env SIZE=$(SIZE) BENCH_REPS=$(BENCH_REPS) $(BIN_SORT_GO)); \
		(env SIZE=$(SIZE) BENCH_REPS=$(BENCH_REPS) $(BIN_SORT_C)); \
		(env SIZE=$(SIZE) BENCH_REPS=$(BENCH_REPS) $(BIN_SORT_RS)); \
		(env BENCH_REPS=$(BENCH_REPS) $(BIN_AOT_SQS) $(SIZE)); \
		(env BENCH_REPS=$(BENCH_REPS) $(BIN_AOT_SMS) $(SIZE)); \
		(env BENCH_REPS=$(BENCH_REPS) $(BIN_AOT_ITER) 45); \
		(env BENCH_REPS=$(BENCH_REPS) $(BIN_AOT_REC) 35); \
		(env $(BIN_FIB_ITER_C) 45); \
		(env $(BIN_FIB_REC_C) 35); \
		(env $(BIN_FIB_ITER_RS) 45); \
		(env $(BIN_FIB_REC_RS) 35); \
	} | awk 'BEGIN{FS="[ =]"} \
	/REPORT/ { \
	  for(i=1;i<=NF;i++){ \
	    if($$i ~ /^impl=/) impl=substr($$i,6); \
	    if($$i ~ /^task=/) task=substr($$i,6); \
	    if($$i ~ /^n=/) n=substr($$i,3); \
	    if($$i ~ /^reps=/) reps=substr($$i,6); \
	    if($$i ~ /^time_ns_avg=/) t=substr($$i,13); \
	    if($$i ~ /^sum=/) chk=substr($$i,5); \
	    if($$i ~ /^result=/) chk=substr($$i,8); \
	  } \
	  if(chk=="") chk="-"; \
	  print impl","task","n","reps","t","chk; next } \
	/TIME_NS:/ { tn=$$2+0; print "legacy,unknown,0,1,"tn",""-"; next }' >> $(CSV_PATH)
	@echo "[csv] $(CSV_PATH)"

$(RESULT_DIR):
	@mkdir -p $(RESULT_DIR)

# ---------- Plotting (gnuplot) ----------

PLOTS_DIR := $(RESULT_DIR)/plots
LATEST_CSV := $(shell ls -1t $(RESULT_DIR)/suite_*.csv 2>/dev/null | head -1)

# Options:
#   PLOT_LOG=1  -> log-scale Y
#   PLOT_REL=1  -> normalize to best (value / min) per chart
.PHONY: plot_csv
plot_csv: | $(PLOTS_DIR)
	@if [ -z "$(LATEST_CSV)" ]; then echo "No CSV found in $(RESULT_DIR). Run 'make bench_all' first."; exit 1; fi
	@if [ ! -s "$(LATEST_CSV)" ]; then echo "CSV is empty: $(LATEST_CSV). Re-run 'make bench_all'."; exit 1; fi
	@echo "[plot] using CSV: $(LATEST_CSV)"

	# ---- Build raw datasets -------------------------------------------------
	@awk -F, 'NR>1 { \
	  t=$$2; sub(/\r$$/,"",t); \
	  if (t=="sort" || t=="sort-qsort" || t=="sort-msort" || t ~ /^sort/) { \
	    label=$$1; \
	    if (label=="tengri-aot") { \
	      if (t=="sort-qsort") label="tengri-aot-qsort"; \
	      else if (t=="sort-msort") label="tengri-aot-msort"; \
	    } \
	    printf "%s\t%s\n", label, $$5; \
	  } \
	}' "$(LATEST_CSV)" > $(PLOTS_DIR)/_sort.raw

	@awk -F, 'NR>1 { t=$$2; sub(/\r$$/,"",t); if (t=="fib_iter") printf "%s\t%s\n", $$1, $$5; }' "$(LATEST_CSV)" > $(PLOTS_DIR)/_fib_iter.raw
	@awk -F, 'NR>1 { t=$$2; sub(/\r$$/,"",t); if (t=="fib_rec")  printf "%s\t%s\n", $$1, $$5; }' "$(LATEST_CSV)" > $(PLOTS_DIR)/_fib_rec.raw

	@if [ ! -s "$(PLOTS_DIR)/_sort.raw" ]; then echo "Empty _sort.raw — CSV probably lacks sort rows."; exit 1; fi
	@if [ ! -s "$(PLOTS_DIR)/_fib_iter.raw" ]; then echo "Empty _fib_iter.raw — CSV probably lacks fib_iter rows."; exit 1; fi
	@if [ ! -s "$(PLOTS_DIR)/_fib_rec.raw" ]; then echo "Empty _fib_rec.raw — CSV probably lacks fib_rec rows."; exit 1; fi

	# ---- Normalize if requested (value := value / min) ----------------------
	@if [ "$(PLOT_REL)" = "1" ]; then \
	  m=$$(awk 'NR==1{min=$$2} $$2<min{min=$$2} END{print min}' $(PLOTS_DIR)/_sort.raw); \
	  awk -v m="$$m" '{printf "%s\t%.6f\n", $$1, $$2/m}' $(PLOTS_DIR)/_sort.raw > $(PLOTS_DIR)/sort.dat; \
	else \
	  cp $(PLOTS_DIR)/_sort.raw $(PLOTS_DIR)/sort.dat; \
	fi

	@if [ "$(PLOT_REL)" = "1" ]; then \
	  m=$$(awk 'NR==1{min=$$2} $$2<min{min=$$2} END{print min}' $(PLOTS_DIR)/_fib_iter.raw); \
	  awk -v m="$$m" '{printf "%s\t%.6f\n", $$1, $$2/m}' $(PLOTS_DIR)/_fib_iter.raw > $(PLOTS_DIR)/fib_iter.dat; \
	else \
	  cp $(PLOTS_DIR)/_fib_iter.raw $(PLOTS_DIR)/fib_iter.dat; \
	fi

	@if [ "$(PLOT_REL)" = "1" ]; then \
	  m=$$(awk 'NR==1{min=$$2} $$2<min{min=$$2} END{print min}' $(PLOTS_DIR)/_fib_rec.raw); \
	  awk -v m="$$m" '{printf "%s\t%.6f\n", $$1, $$2/m}' $(PLOTS_DIR)/_fib_rec.raw > $(PLOTS_DIR)/fib_rec.dat; \
	else \
	  cp $(PLOTS_DIR)/_fib_rec.raw $(PLOTS_DIR)/fib_rec.dat; \
	fi

	# ---- Gnuplot scripts ----------------------------------------------------
	@echo 'set terminal pngcairo size 1000,600'                                    >  $(BIN_DIR)/plot_common.gp
	@echo 'set datafile separator "\t"'                                           >> $(BIN_DIR)/plot_common.gp
	@echo 'set style data histograms'                                             >> $(BIN_DIR)/plot_common.gp
	@echo 'set style fill solid 1.00 border -1'                                   >> $(BIN_DIR)/plot_common.gp
	@echo 'set boxwidth 0.8'                                                      >> $(BIN_DIR)/plot_common.gp
	@echo 'set grid ytics'                                                        >> $(BIN_DIR)/plot_common.gp
	@if [ "$(PLOT_LOG)" = "1" ]; then echo 'set logscale y'; else echo 'unset logscale y'; fi >> $(BIN_DIR)/plot_common.gp

	@YLABEL=$$( [ "$(PLOT_REL)" = "1" ] && echo 'factor of best (×)' || echo 'time_ns_avg (ns)' ); \
	TITLE_S=$$( [ "$(PLOT_REL)" = "1" ] && echo 'Sort benchmark — normalized (smaller is better)' || echo 'Sort benchmark (smaller is better)' ); \
	TITLE_I=$$( [ "$(PLOT_REL)" = "1" ] && echo 'fib_iter — normalized (smaller is better)' || echo 'fib_iter benchmark (smaller is better)' ); \
	TITLE_R=$$( [ "$(PLOT_REL)" = "1" ] && echo 'fib_rec — normalized (smaller is better)'  || echo 'fib_rec benchmark (smaller is better)' ); \
	{ \
	  cat $(BIN_DIR)/plot_common.gp; \
	  echo "set output \"benchmarks/results/plots/sort.png\""; \
	  echo "set ylabel \"$$YLABEL\""; \
	  echo "set title \"$$TITLE_S\""; \
	  echo "plot \"benchmarks/results/plots/sort.dat\" using 2:xtic(1) title \"\""; \
	} > $(BIN_DIR)/plot_sort.gp; \
	gnuplot $(BIN_DIR)/plot_sort.gp; \
	{ \
	  cat $(BIN_DIR)/plot_common.gp; \
	  echo "set output \"benchmarks/results/plots/fib_iter.png\""; \
	  echo "set ylabel \"$$YLABEL\""; \
	  echo "set title \"$$TITLE_I\""; \
	  echo "plot \"benchmarks/results/plots/fib_iter.dat\" using 2:xtic(1) title \"\""; \
	} > $(BIN_DIR)/plot_fib_iter.gp; \
	gnuplot $(BIN_DIR)/plot_fib_iter.gp; \
	{ \
	  cat $(BIN_DIR)/plot_common.gp; \
	  echo "set output \"benchmarks/results/plots/fib_rec.png\""; \
	  echo "set ylabel \"$$YLABEL\""; \
	  echo "set title \"$$TITLE_R\""; \
	  echo "plot \"benchmarks/results/plots/fib_rec.dat\" using 2:xtic(1) title \"\""; \
	} > $(BIN_DIR)/plot_fib_rec.gp; \
	gnuplot $(BIN_DIR)/plot_fib_rec.gp

	@echo "[plots] benchmarks/results/plots/sort.png"
	@echo "[plots] benchmarks/results/plots/fib_iter.png"
	@echo "[plots] benchmarks/results/plots/fib_rec.png"

$(PLOTS_DIR):
	@mkdir -p $(PLOTS_DIR)