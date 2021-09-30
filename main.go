package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for k, vals := range r.Header {
			for _, v := range vals {
				w.Header().Add(k, v)
			}

		}

		if envVar, has := os.LookupEnv("VERSION"); has {
			w.Header().Add("ENV_VAR", envVar)
		}

		w.WriteHeader(200)
		w.Write([]byte("OK"))

		log.Printf("%s:%d", getIP(r), 200)

	})
	
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	http.ListenAndServe(":8888", nil)

}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}