package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", factorRandom)
	http.ListenAndServe(":8080", nil)
}

func factorRandom(w http.ResponseWriter, r *http.Request) {
	nr := int64(rand.Intn(100000000000000) + 100000000000000)
	nr_str := strconv.FormatInt(nr, 10)
	if nr < 1 {
		fmt.Fprintf(w, "\nFactors of", nr_str, "not computed")
		return
	}
	//fmt.Fprintf(w, "\nFactors of %d: ", nr_str)
	fs := make([]int64, 1)
	fs[0] = 1
	apf := func(p int64, e int) {
		n := len(fs)
		for i, pp := 0, p; i < e; i, pp = i+1, pp*p {
			for j := 0; j < n; j++ {
				fs = append(fs, fs[j]*pp)
			}
		}
	}
	e := 0
	for ; nr&1 == 0; e++ {
		nr >>= 1
	}
	apf(2, e)
	for d := int64(3); nr > 1; d += 2 {
		if d*d > nr {
			d = nr
		}
		for e = 0; nr%d == 0; e++ {
			nr /= d
		}
		if e > 0 {
			apf(d, e)
		}
	}

	fmt.Fprintf(w, "Number of factors for %s is = %d", nr_str, len(fs))
}
