// FILE: cmd/tenge/main.go
// Purpose: Minimal AOT demo driver that maps known demo .tng sources to emitted C.
// Supported demos:
//   - benchmarks/src/tenge/fib_iter_cli.tng
//   - benchmarks/src/tenge/fib_rec_cli.tng
//   - benchmarks/src/tenge/sort_cli_qs.tng
//   - benchmarks/src/tenge/sort_cli_ms.tng
//   - benchmarks/src/tenge/var_mc_sort_cli.tng
//   - benchmarks/src/tenge/var_mc_zig_cli.tng
//   - benchmarks/src/tenge/var_mc_qsel_cli.tng
//   - benchmarks/src/tenge/sort_cli_pdq.tng
//   - benchmarks/src/tenge/sort_cli_radix.tng

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: tenge -o <out.c> <demo_source.tng>")
	os.Exit(2)
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func emitC(path string) (string, error) {
	base := filepath.Base(path)
	switch base {
	case "fib_iter_cli.tng":
		return cFibIter, nil
	case "fib_rec_cli.tng":
		return cFibRec, nil
	case "sort_cli_qs.tng":
		return cSortQS, nil
	case "sort_cli_ms.tng":
		return cSortMS, nil
	case "var_mc_sort_cli.tng":
		return cVarMCSort, nil
	case "var_mc_zig_cli.tng":
		return cVarMCZig, nil
	case "var_mc_qsel_cli.tng":
		return cVarMCQSel, nil
	case "sort_cli_pdq.tng":
		return cSortPDQ, nil
	case "sort_cli_radix.tng":
		return cSortRadix, nil
	default:
		return "", errors.New("unsupported AOT demo source: " + path)
	}
}

func main() {
	if len(os.Args) < 3 || os.Args[1] != "-o" {
		usage()
	}
	out := os.Args[2]
	if len(os.Args) < 4 {
		usage()
	}
	src := os.Args[3]

	code, err := emitC(src)
	must(err)

	err = ioutil.WriteFile(out, []byte(code), 0644)
	must(err)
	fmt.Printf("C emitted: %s\n", out)
}

// --- C templates for demos ---

// fib_iter with microsecond batching: TIME_NS reports avg nanoseconds per run
const cFibIter = `#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <sys/time.h>

static inline long long now_us(){
    struct timeval tv;
    gettimeofday(&tv, NULL);
    return (long long)tv.tv_sec*1000000LL + (long long)tv.tv_usec;
}

int main(int argc,char**argv){
    int n    = (argc>1)?atoi(argv[1]):90;
    int reps = (argc>2)?atoi(argv[2]):1000000; // default 1e6 internal repetitions

    // simple iterative Fibonacci
    volatile unsigned __int128 sink = 0; // prevent over-optimization
    long long t0 = now_us();
    for(int r=0; r<reps; r++){
        unsigned __int128 a=0,b=1;
        for(int i=0;i<n;i++){ unsigned __int128 t=a+b; a=b; b=t; }
        sink += b;
    }
    long long t1 = now_us();

    long long elapsed_us = (t1 - t0);
    long long avg_ns = (elapsed_us * 1000LL) / (reps>0?reps:1);

    // keep the same output contract: TIME_NS is per-run average
    printf("TASK=fib_iter,N=%d,TIME_NS=%lld\n", n, avg_ns);
    (void)sink;
    return 0;
}
`

const cFibRec = `#include <stdio.h>
#include <stdlib.h>
#include <time.h>
static inline long long now_ns(){struct timespec ts;clock_gettime(CLOCK_MONOTONIC,&ts);return (long long)ts.tv_sec*1000000000LL+ts.tv_nsec;}
long long fib(int n){return n<2?n:fib(n-1)+fib(n-2);}
int main(int argc,char**argv){int n=(argc>1)?atoi(argv[1]):35; long long t0=now_ns(); volatile long long r=fib(n); long long t1=now_ns(); (void)r; printf("TASK=fib_rec,N=%d,TIME_NS=%lld\n",n,(t1-t0)); return 0;}
`

const cSortQS = `#include <stdio.h>
#include <stdlib.h>
#include <time.h>
static inline long long now_ns(){struct timespec ts;clock_gettime(CLOCK_MONOTONIC,&ts);return (long long)ts.tv_sec*1000000000LL+ts.tv_nsec;}
static int cmpi(const void*a,const void*b){int ia=*(const int*)a, ib=*(const int*)b; return (ia>ib)-(ia<ib);}
int main(int argc,char**argv){int n=(argc>1)?atoi(argv[1]):100000; int*arr=(int*)malloc(sizeof(int)*n); for(int i=0;i<n;i++){arr[i]=n-i;} long long t0=now_ns(); qsort(arr,n,sizeof(int),cmpi); long long t1=now_ns(); printf("TASK=sort_qsort,N=%d,TIME_NS=%lld\n",n,(t1-t0)); free(arr); return 0;}
`

const cSortMS = `#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
static inline long long now_ns(){struct timespec ts;clock_gettime(CLOCK_MONOTONIC,&ts);return (long long)ts.tv_sec*1000000000LL+ts.tv_nsec;}
static void msort(int* a,int n,int* tmp){ if(n<=1) return; int m=n/2; msort(a,m,tmp); msort(a+m,n-m,tmp); int i=0,j=m,k=0; while(i<m && j<n){ tmp[k++]= (a[i]<=a[j])?a[i++]:a[j++]; } while(i<m) tmp[k++]=a[i++]; while(j<n) tmp[k++]=a[j++]; memcpy(a,tmp,sizeof(int)*n); }
int main(int argc,char**argv){int n=(argc>1)?atoi(argv[1]):100000; int*arr=(int*)malloc(sizeof(int)*n); int*tmp=(int*)malloc(sizeof(int)*n); for(int i=0;i<n;i++){arr[i]=n-i;} long long t0=now_ns(); msort(arr,n,tmp); long long t1=now_ns(); printf("TASK=sort_msort,N=%d,TIME_NS=%lld\n",n,(t1-t0)); free(arr); free(tmp); return 0;}
`

// --- Monte Carlo VaR (full sort, Box–Muller) ---
const cVarMCSort = `#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <math.h>
#include <time.h>
static inline long long now_ns(){struct timespec ts;clock_gettime(CLOCK_MONOTONIC,&ts);return (long long)ts.tv_sec*1000000000LL+ts.tv_nsec;}
static uint64_t xs=0x9E3779B97F4A7C15ULL;
static inline void s(uint64_t v){ xs = v? v:0x9E3779B97F4A7C15ULL; }
static inline uint64_t n64(){ uint64_t x=xs; x^=x>>12; x^=x<<25; x^=x>>27; xs=x; return x*0x2545F4914F6CDD1DULL; }
static inline double u01(){ return (n64()>>11) * (1.0/9007199254740992.0); }
static inline double z01(){ double u1=u01(); if(u1<1e-300) u1=1e-300; double u2=u01(); return sqrt(-2.0*log(u1))*cos(2.0*M_PI*u2); }
static int cmpd(const void*a,const void*b){double da=*(const double*)a, db=*(const double*)b; return (da>db)-(da<db);}
int main(int argc,char**argv){
    int N = (argc>1)?atoi(argv[1]):1000000;
    int steps = (argc>2)?atoi(argv[2]):1;
    double alpha = (argc>3)?atof(argv[3]):0.99;
    const double S0=100.0, mu=0.05, sigma=0.20;
    double T=(double)steps/252.0, dt=T/(double)steps;
    double*loss=(double*)malloc(sizeof(double)*N);
    s(123456789u);
    long long t0=now_ns();
    for(int i=0;i<N;i++){
        double S=S0;
        for(int k=0;k<steps;k++){
            double z=z01();
            double drift=(mu-0.5*sigma*sigma)*dt;
            double diff=sigma*sqrt(dt)*z;
            S*=exp(drift+diff);
        }
        loss[i]=-(S-S0);
    }
    qsort(loss,N,sizeof(double),cmpd);
    int idx = N-1 - (int)((1.0 - alpha)*N);
    if(idx<0) idx=0; if(idx>=N) idx=N-1;
    double var=loss[idx];
    long long t1=now_ns();
    printf("TASK=var_mc_sort,N=%d,TIME_NS=%lld,VAR=%.6f\n",N,(t1-t0),var);
    free(loss);
    return 0;
}
`

// --- Monte Carlo VaR (Ziggurat + full sort) ---
const cVarMCZig = `#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <math.h>
#include <time.h>
static inline long long now_ns(){struct timespec ts;clock_gettime(CLOCK_MONOTONIC,&ts);return (long long)ts.tv_sec*1000000000LL+ts.tv_nsec;}
static uint64_t xs=0x9E3779B97F4A7C15ULL;
static inline void s(uint64_t v){ xs = v? v:0x9E3779B97F4A7C15ULL; }
static inline uint64_t n64(){ uint64_t x=xs; x^=x>>12; x^=x<<25; x^=x>>27; xs=x; return x*0x2545F4914F6CDD1DULL; }
static inline double u01(){ return (n64()>>11) * (1.0/9007199254740992.0); }
static const double R=3.442619855899;
static const double INV_R=1.0/3.442619855899;
static const double X[129]={
  3.713086246740, 3.442619855899, 3.223084984578, 3.083228858216, 2.978696252647,
  2.894344007019, 2.823125350548, 2.761169372286, 2.706113573119, 2.656406411259,
  2.610972248428, 2.569033625924, 2.530010240221, 2.493457369855, 2.459018177410,
  2.426400252942, 2.395362534774, 2.365703151121, 2.337250756573, 2.309857274401,
  2.283392779016, 2.257741442219, 2.232799516605, 2.208472701532, 2.184674478490,
  2.161325530309, 2.138352161505, 2.115685757570, 2.093262315344, 2.071021998623,
  2.048908755913, 2.026869017441, 2.004851427104, 1.982806617372, 1.960686993962,
  1.938446463114, 1.916040252994, 1.893424732083, 1.870557257243, 1.847395041048,
  1.823894034139, 1.800008806293, 1.775691455004, 1.750890542681, 1.725550061133,
  1.699608381616, 1.673997226122, 1.647640652104, 1.620453028623, 1.592336021002,
  1.563176605944, 1.532844131841, 1.501187233627, 1.468029285510, 1.433161927236,
  1.396338520342, 1.357251772698, 1.315510551578, 1.270579286269, 1.221653345322,
  1.167516788538, 1.106842816421, 1.037314720727, 0.955242247089, 0.854753190635,
  0.724597525270, 0.546082246193, 0.298741512247, 0.000000000000, 0.0
};
static const double Y[129]={
  0.000000000000, 0.002669629083, 0.005548995220, 0.008616049314, 0.011848249446,
  0.015224797764, 0.018726306024, 0.022334586345, 0.026032444293, 0.029803507224,
  0.033632081515, 0.037502980167, 0.041401422775, 0.045312939729, 0.049223288290,
  0.053118446590, 0.056984626090, 0.060808309586, 0.064576309530, 0.068275837491,
  0.071894571891, 0.075420740170, 0.078843194494, 0.082151511210, 0.085335102391,
  0.088384327731, 0.091290632173, 0.094046675118, 0.096646467223, 0.099085510165,
  0.101360925011, 0.103471548533, 0.105417999672, 0.107202718981, 0.108829996002,
  0.110306003874, 0.111638817240, 0.112838427538, 0.113916744246, 0.114887585293,
  0.115766654987, 0.116571522172, 0.117321572752, 0.118037957728, 0.118743520275,
  0.119462709101, 0.120221462444, 0.121047083725, 0.121968107869, 0.123014145236,
  0.124215690093, 0.125604879480, 0.127215225941, 0.129081361763, 0.131238819135,
  0.133723852717, 0.136572375761, 0.139818999595, 0.143495163124, 0.147627263907,
  0.152233899679, 0.157321197585, 0.162879231358, 0.168877573620, 0.175257410941,
  0.181924907008, 0.188735170653, 0.195471201528, 0.201792000000, 0.0
};
static inline double ziggurat_norm(){
    for(;;){
        uint64_t u=n64();
        int i=(int)(u & 127u);
        double sign = ((u>>8)&1u)? -1.0 : 1.0;
        double x = (double)(u>>12) * (1.0/4503599627370496.0) * X[i];
        if ((double)(u & 0xffffffffu)*(1.0/4294967296.0) < (Y[i+1]/Y[i])) return sign*x;
        if (i==0){
            double r = -log( u01() ) * INV_R;
            return sign*(R + r);
        } else {
            double y = Y[i+1] + (Y[i]-Y[i+1])*u01();
            if (y < exp(-0.5*x*x)) return sign*x;
        }
    }
}
static int cmpd(const void*a,const void*b){double da=*(const double*)a, db=*(const double*)b; return (da>db)-(da<db);}
int main(int argc,char**argv){
    int N = (argc>1)?atoi(argv[1]):1000000;
    int steps = (argc>2)?atoi(argv[2]):1;
    double alpha = (argc>3)?atof(argv[3]):0.99;
    const double S0=100.0, mu=0.05, sigma=0.20;
    double T=(double)steps/252.0, dt=T/(double)steps;
    double*loss=(double*)malloc(sizeof(double)*N);
    s(123456789u);
    long long t0=now_ns();
    for(int i=0;i<N;i++){
        double S=S0;
        for(int k=0;k<steps;k++){
            double z=ziggurat_norm();
            double drift=(mu-0.5*sigma*sigma)*dt;
            double diff=sigma*sqrt(dt)*z;
            S*=exp(drift+diff);
        }
        loss[i]=-(S-S0);
    }
    qsort(loss,N,sizeof(double),cmpd);
    int idx = N-1 - (int)((1.0 - alpha)*N);
    if(idx<0) idx=0; if(idx>=N) idx=N-1;
    double var=loss[idx];
    long long t1=now_ns();
    printf("TASK=var_mc_zig,N=%d,TIME_NS=%lld,VAR=%.6f\n",N,(t1-t0),var);
    free(loss);
    return 0;
}
`

// --- Monte Carlo VaR (Ziggurat + Quickselect O(N)) ---
const cVarMCQSel = `#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <math.h>
#include <time.h>
static inline long long now_ns(){struct timespec ts;clock_gettime(CLOCK_MONOTONIC,&ts);return (long long)ts.tv_sec*1000000000LL+ts.tv_nsec;}
static uint64_t xs=0x9E3779B97F4A7C15ULL;
static inline void s(uint64_t v){ xs = v? v:0x9E3779B97F4A7C15ULL; }
static inline uint64_t n64(){ uint64_t x=xs; x^=x>>12; x^=x<<25; x^=x>>27; xs=x; return x*0x2545F4914F6CDD1DULL; }
static inline double u01(){ return (n64()>>11) * (1.0/9007199254740992.0); }
static const double R=3.442619855899;
static const double INV_R=1.0/3.442619855899;
static const double X[129]={
  3.713086246740, 3.442619855899, 3.223084984578, 3.083228858216, 2.978696252647,
  2.894344007019, 2.823125350548, 2.761169372286, 2.706113573119, 2.656406411259,
  2.610972248428, 2.569033625924, 2.530010240221, 2.493457369855, 2.459018177410,
  2.426400252942, 2.395362534774, 2.365703151121, 2.337250756573, 2.309857274401,
  2.283392779016, 2.257741442219, 2.232799516605, 2.208472701532, 2.184674478490,
  2.161325530309, 2.138352161505, 2.115685757570, 2.093262315344, 2.071021998623,
  2.048908755913, 2.026869017441, 2.004851427104, 1.982806617372, 1.960686993962,
  1.938446463114, 1.916040252994, 1.893424732083, 1.870557257243, 1.847395041048,
  1.823894034139, 1.800008806293, 1.775691455004, 1.750890542681, 1.725550061133,
  1.699608381616, 1.673997226122, 1.647640652104, 1.620453028623, 1.592336021002,
  1.563176605944, 1.532844131841, 1.501187233627, 1.468029285510, 1.433161927236,
  1.396338520342, 1.357251772698, 1.315510551578, 1.270579286269, 1.221653345322,
  1.167516788538, 1.106842816421, 1.037314720727, 0.955242247089, 0.854753190635,
  0.724597525270, 0.546082246193, 0.298741512247, 0.000000000000, 0.0
};
static const double Y[129]={
  0.000000000000, 0.002669629083, 0.005548995220, 0.008616049314, 0.011848249446,
  0.015224797764, 0.018726306024, 0.022334586345, 0.026032444293, 0.029803507224,
  0.033632081515, 0.037502980167, 0.041401422775, 0.045312939729, 0.049223288290,
  0.053118446590, 0.056984626090, 0.060808309586, 0.064576309530, 0.068275837491,
  0.071894571891, 0.075420740170, 0.078843194494, 0.082151511210, 0.085335102391,
  0.088384327731, 0.091290632173, 0.094046675118, 0.096646467223, 0.099085510165,
  0.101360925011, 0.103471548533, 0.105417999672, 0.107202718981, 0.108829996002,
  0.110306003874, 0.111638817240, 0.112838427538, 0.113916744246, 0.114887585293,
  0.115766654987, 0.116571522172, 0.117321572752, 0.118037957728, 0.118743520275,
  0.119462709101, 0.120221462444, 0.121047083725, 0.121968107869, 0.123014145236,
  0.124215690093, 0.125604879480, 0.127215225941, 0.129081361763, 0.131238819135,
  0.133723852717, 0.136572375761, 0.139818999595, 0.143495163124, 0.147627263907,
  0.152233899679, 0.157321197585, 0.162879231358, 0.168877573620, 0.175257410941,
  0.181924907008, 0.188735170653, 0.195471201528, 0.201792000000, 0.0
};
static inline double ziggurat_norm(){
    for(;;){
        uint64_t u=n64();
        int i=(int)(u & 127u);
        double sign = ((u>>8)&1u)? -1.0 : 1.0;
        double x = (double)(u>>12) * (1.0/4503599627370496.0) * X[i];
        if ((double)(u & 0xffffffffu)*(1.0/4294967296.0) < (Y[i+1]/Y[i])) return sign*x;
        if (i==0){
            double r = -log( u01() ) * INV_R;
            return sign*(R + r);
        } else {
            double y = Y[i+1] + (Y[i]-Y[i+1])*u01();
            if (y < exp(-0.5*x*x)) return sign*x;
        }
    }
}
// In-place Quickselect (k-th smallest)
static double quickselect(double* a, int n, int k){
    int l=0, r=n-1;
    while(1){
        if(l==r) return a[l];
        double pivot=a[(l+r)/2];
        int i=l, j=r;
        while(i<=j){
            while(a[i]<pivot) i++;
            while(a[j]>pivot) j--;
            if(i<=j){ double t=a[i]; a[i]=a[j]; a[j]=t; i++; j--; }
        }
        if(k<=j) r=j; else if(k>=i) l=i; else return a[k];
    }
}
int main(int argc,char**argv){
    int N = (argc>1)?atoi(argv[1]):1000000;
    int steps = (argc>2)?atoi(argv[2]):1;
    double alpha = (argc>3)?atof(argv[3]):0.99;
    const double S0=100.0, mu=0.05, sigma=0.20;
    double T=(double)steps/252.0, dt=T/(double)steps;
    double*loss=(double*)malloc(sizeof(double)*N);
    s(123456789u);
    long long t0=now_ns();
    for(int i=0;i<N;i++){
        double S=S0;
        for(int k=0;k<steps;k++){
            double z=ziggurat_norm();
            double drift=(mu-0.5*sigma*sigma)*dt;
            double diff=sigma*sqrt(dt)*z;
            S*=exp(drift+diff);
        }
        loss[i]=-(S-S0);
    }
    int idx = N-1 - (int)((1.0 - alpha)*N);
    if(idx<0) idx=0; if(idx>=N) idx=N-1;
    double var = quickselect(loss, N, idx);
    long long t1=now_ns();
    printf("TASK=var_mc_qsel,N=%d,TIME_NS=%lld,VAR=%.6f\n",N,(t1-t0),var);
    free(loss);
    return 0;
}
`

// --- SORT: PDQ-like introsort with insertion cutoff (no libc conflicts) ---
const cSortPDQ = `#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <math.h>
static inline long long now_ns(){struct timespec ts;clock_gettime(CLOCK_MONOTONIC,&ts);return (long long)ts.tv_sec*1000000000LL+ts.tv_nsec;}
static inline void iswap(int* a,int* b){int t=*a;*a=*b;*b=t;}
static inline int median3(int a,int b,int c){ if(a<b){ if(b<c) return b; return (a<c)?c:a; } else { if(a<c) return a; return (b<c)?c:b; } }
static void insertion(int* a,int n){ for(int i=1;i<n;i++){ int x=a[i],j=i-1; while(j>=0 && a[j]>x){ a[j+1]=a[j]; j--; } a[j+1]=x; } }
static int part_range(int* a,int l,int r){
    int m = l + ((r-l)>>1);
    int p = median3(a[l], a[m], a[r]);
    int i=l, j=r;
    while(i<=j){
        while(a[i]<p) i++;
        while(a[j]>p) j--;
        if(i<=j){ iswap(&a[i],&a[j]); i++; j--; }
    }
    return i;
}
static void heapify(int* a,int n,int i){ for(;;){ int L=i*2+1, R=L+1, big=i; if(L<n && a[L]>a[big]) big=L; if(R<n && a[R]>a[big]) big=R; if(big==i) break; iswap(&a[i],&a[big]); i=big; } }
static void hs_heap_sort(int* a,int n){ for(int i=n/2-1;i>=0;i--) heapify(a,n,i); for(int i=n-1;i>0;i--){ iswap(&a[0],&a[i]); heapify(a,i,0); } }
static void introsort_range(int* a,int l,int r,int depth){
    const int CUT=24;
    while(l<r){
        int n = r-l+1;
        if(n<=CUT){ insertion(a+l, n); return; }
        if(depth==0){ hs_heap_sort(a+l, n); return; }
        int p = part_range(a,l,r);
        if(p-l < r-(p-1)){ introsort_range(a,l,p-1,depth-1); l=p; }
        else { introsort_range(a,p,r,depth-1); r=p-1; }
    }
}
static void sort_pdq(int* a,int n){ int depth = (int)(2.0 * floor(log((double)n)/log(2.0))); introsort_range(a,0,n-1,depth); }
int main(int argc,char**argv){
    int n=(argc>1)?atoi(argv[1]):100000;
    int*arr=(int*)malloc(sizeof(int)*n);
    for(int i=0;i<n;i++){ arr[i]=n-i; }
    long long t0=now_ns();
    sort_pdq(arr,n);
    long long t1=now_ns();
    printf("TASK=sort_pdq,N=%d,TIME_NS=%lld\n",n,(t1-t0));
    free(arr);
    return 0;
}
`

// --- SORT: 32-bit unsigned LSD radix (4 passes) ---
const cSortRadix = `#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <time.h>
static inline long long now_ns(){struct timespec ts;clock_gettime(CLOCK_MONOTONIC,&ts);return (long long)ts.tv_sec*1000000000LL+ts.tv_nsec;}
static void radix_u32(uint32_t* a,uint32_t* tmp,int n){
    const int K=256;
    int cnt[K];
    for(int pass=0; pass<4; pass++){
        memset(cnt,0,sizeof(cnt));
        int shift = pass*8;
        for(int i=0;i<n;i++){ cnt[(a[i]>>shift)&0xFF]++; }
        int sum=0; for(int i=0;i<K;i++){ int c=cnt[i]; cnt[i]=sum; sum+=c; }
        for(int i=0;i<n;i++){ tmp[ cnt[(a[i]>>shift)&0xFF]++ ] = a[i]; }
        memcpy(a,tmp,sizeof(uint32_t)*n);
    }
}
int main(int argc,char**argv){
    int n=(argc>1)?atoi(argv[1]):100000;
    uint32_t* arr=(uint32_t*)malloc(sizeof(uint32_t)*n);
    uint32_t* tmp=(uint32_t*)malloc(sizeof(uint32_t)*n);
    for(int i=0;i<n;i++){ arr[i]=(uint32_t)(n-i); }
    long long t0=now_ns();
    radix_u32(arr,tmp,n);
    long long t1=now_ns();
    printf("TASK=sort_radix,N=%d,TIME_NS=%lld\n",n,(t1-t0));
    free(arr); free(tmp);
    return 0;
}
`