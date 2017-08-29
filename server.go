package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

//func Run(addr string, sslAddr string, ssl map[string]string, insecureHandler, secureHandler http.Handler) chan error {
//
//	errs := make(chan error)
//
//	// Starting HTTP server
//	go func() {
//		log.Printf("Staring HTTP service on %s ...", addr)
//
//		if err := http.ListenAndServe(addr, insecureHandler); err != nil {
//			errs <- err
//		}
//
//	}()
//
//	// Starting HTTPS server
//	go func() {
//		log.Printf("Staring HTTPS service on %s ...", sslAddr)
//		if err := http.ListenAndServeTLS(sslAddr, ssl["cert"], ssl["key"], secureHandler); err != nil {
//			errs <- err
//		}
//	}()
//
//	return errs
//}
//
//func insecureHandler(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "text/plain")
//	w.Write([]byte("Yay!! insecure server.\n"))
//}

func secureHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Yay!! secure server.\n"))
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Yay!! secure server.\n"))
}

func main() {
	s := mux.NewRouter()
	s.HandleFunc("/messages", addHandler).Methods("POST")
	s.HandleFunc("/messages", secureHandler).Methods("GET")

	log.Printf("Staring HTTPS service on %s ...", ":443")
	if err := http.ListenAndServeTLS(":443", "localhost.crt", "server.key", s); err != nil {
		panic(err)
	}
}
