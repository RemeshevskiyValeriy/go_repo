package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"example.com/pprof-lab/internal/work"
)

func main() {
	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		n := 38
		defer work.TimeIt("Fib(38)")()
		res := work.FibFast(n)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(fmt.Sprintf("%d\n", res)))
	})

	log.Println("Server on :8080; pprof on /debug/pprof/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
