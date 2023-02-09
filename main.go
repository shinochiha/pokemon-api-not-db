package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pokemon-api/handlers"
)

func main() {
	r := mux.NewRouter()
	r.Use(commonMiddleware)
	r.HandleFunc("/pokemon", handlers.GetPokemonListHandler).Methods("GET")
	r.HandleFunc("/pokemon/{id}", handlers.GetPokemonByIDHandler).Methods("GET")
	r.HandleFunc("/pokemon/{id}", handlers.CatchPokemonByIdHandler).Methods("POST")
	r.HandleFunc("/mypokemon", handlers.GetMyPokemonHandler).Methods("GET")
	r.HandleFunc("/mypokemon/{id}", handlers.ReleasePokemonHandler).Methods("DELETE")
	r.HandleFunc("/mypokemon/{id}/{nickname}", handlers.ChangeNicknameHandler).Methods("PUT")
	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
