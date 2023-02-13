package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pokemon-api/models"
)

func GetPokemonByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid pokemon ID", http.StatusBadRequest)
		return
	}
	pokemon, err := models.GetPokemonByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(pokemon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetPokemonListHandler(w http.ResponseWriter, r *http.Request) {
	offsetStr := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	limitPokeApiStr := r.URL.Query().Get("limit_poke_api")
	limitPokeApi, err := strconv.Atoi(limitPokeApiStr)
	if err != nil {
		limitPokeApi = 20
	}

	pokemons, err := models.GetPokemonList(limitPokeApi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pokemonList := &models.PokemonList{
		Count:    len(pokemons.Results),
		Next:     "",
		Previous: "",
		Results:  pokemons.Results[offset : offset+limit],
	}

	if offset+limit < len(pokemons.Results) {
		pokemonList.Next = fmt.Sprintf("/pokemon?offset=%d&limit=%d", offset+limit, limit)
	}
	if offset > 0 {
		pokemonList.Previous = fmt.Sprintf("/pokemon?offset=%d&limit=%d", offset-limit, limit)
	}

	err = json.NewEncoder(w).Encode(pokemonList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CatchPokemonByIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	pokemon, err := models.CatchPokemonById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(pokemon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetMyPokemonHandler(w http.ResponseWriter, r *http.Request) {
	count, myPokemon := models.GetMyPokemon()
	response := map[string]interface{}{
		"count":   count,
		"results": myPokemon,
	}
	json.NewEncoder(w).Encode(response)
}

func ReleasePokemonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid value for id", http.StatusBadRequest)
		return
	}

	queryValues := r.URL.Query()
	primeNumber, err := strconv.Atoi(queryValues.Get("release"))
	if err != nil {
		http.Error(w, "Invalid value for release", http.StatusBadRequest)
		return
	}

	released := models.ReleasePokemon(id, primeNumber)
	if released {
		w.Write([]byte(`{"message": "Pokemon has been released successfully."}`))
	} else {
		http.Error(w, "Can't release the pokemon with this number", http.StatusBadRequest)
	}
}

func ChangeNicknameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid value for id", http.StatusBadRequest)
		return
	}
	nickname := vars["nickname"]
	if nickname == "" {
		http.Error(w, "Nickname is required", http.StatusBadRequest)
		return
	}

	err = models.ChangeNickname(id, nickname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	newNickname := fmt.Sprintf("%s-%d", nickname, fibonacci())
	w.Write([]byte(fmt.Sprintf(`{"message": "Nickname changed successfully", "nickname": "%s"}`, newNickname)))
}
