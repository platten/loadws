package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Globals
var MAX_SIZE int

func main() {
	var port int

	flag.IntVar(&port, "port", 8080, "port")
	flag.IntVar(&MAX_SIZE, "size", 1000000000000000, "order of magnitude for factorization")
	flag.Parse()

	wsPort := ":" + strconv.Itoa(port)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", factorRandom)
	fmt.Println("Running server on port: ", port)
	log.Fatal(http.ListenAndServe(wsPort, nil))
}

func factorRandom(w http.ResponseWriter, r *http.Request) {
	nr := int64(rand.Intn(MAX_SIZE) + MAX_SIZE)
	nrStr := strconv.FormatInt(nr, 10)

	// Taken from RosettaCode: https://rosettacode.org/wiki/Factors_of_an_integer#Go
	if nr < 1 {
		fmt.Fprintf(w, "\nFactors of %s not computed", nrStr)
		return
	}
	//fmt.Fprintf(w, "\nFactors of %d: ", nrStr)
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

	fmt.Fprintf(w, "Number of factors for %s is = %d", nrStr, len(fs))
}
