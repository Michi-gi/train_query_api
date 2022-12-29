package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func main() {
	router := chi.NewRouter()

	router.Route("/station/{stationId}", func(r chi.Router) {
		r.Get("/", StationController)
		r.Route("/{tableId}", func(r chi.Router) {
			r.Get("/", TableController)
			r.Get("/{trainId}", TrainController)
		})
	})
	router.Get("/search", SearchController)

	port := os.Getenv("PORT")
	fmt.Print("port: " + port)
	if err := http.ListenAndServe(":" + port, router); err != nil {
		fmt.Print(err)
	}
}
