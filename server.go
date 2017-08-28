package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Run(addr string, sslAddr string, ssl map[string]string, insecureHandler, secureHandler http.Handler) chan error {

	errs := make(chan error)

	// Starting HTTP server
	go func() {
		log.Printf("Staring HTTP service on %s ...", addr)

		if err := http.ListenAndServe(addr, insecureHandler); err != nil {
			errs <- err
		}

	}()

	// Starting HTTPS server
	go func() {
		log.Printf("Staring HTTPS service on %s ...", sslAddr)
		if err := http.ListenAndServeTLS(sslAddr, ssl["cert"], ssl["key"], secureHandler); err != nil {
			errs <- err
		}
	}()

	return errs
}

//func sampleHandler(w http.ResponseWriter, req *http.Request) {
//	w.Header().Set("Content-Type", "text/plain")
//	w.Write([]byte("This is an example server.\n"))
//}

func insecureHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Yay!! insecure server.\n"))
}

func secureHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Yay!! secure server.\n"))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Yay!! default server.\n"))
}

func main() {
	i := mux.NewRouter()
	s := mux.NewRouter()
	i.HandleFunc("/", defaultHandler)
	//r.HandleFunc("/products", handler).Methods("POST")
	i.HandleFunc("/i", insecureHandler).Methods("GET")
	s.HandleFunc("/s", secureHandler).Methods("GET")
	//http.Handle("/", r)
	//http.HandleFunc("/", sampleHandler)

	errs := Run(":8081", ":443", map[string]string{
		"cert": "localhost.crt",
		"key":  "server.key",
	},
		i,
		s)

	// This will run forever until channel receives error
	select {
	case err := <-errs:
		log.Printf("Could not start serving service due to (error: %s)", err)
	}

}
