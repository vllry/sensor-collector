package webview

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

type view struct {
	temperatureData map[int32]float32
}

func (v *view) index(w http.ResponseWriter, r *http.Request) {
	var s string
	for key, value := range v.temperatureData {
		s += fmt.Sprintf("Sensor %v: %v\n", key, value)
	}

	w.Write([]byte(s)) // Ugh
}

func RunWebserver(temperatureData map[int32]float32) {
	viewStruct := view{
		temperatureData: temperatureData,
	}

	router := mux.NewRouter()

	router.HandleFunc("/", viewStruct.index).Methods("GET")

	router.Use(loggingMiddleware)
	fmt.Println("Listening...")
	log.Fatal(http.ListenAndServe(":80", router))
}
