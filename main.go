package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Struct untuk Detail Pokemon
type PokemonDetail struct {
	ID      int                    `json:"id"`
	Name    string                 `json:"name"`
	Image   map[string]interface{} `json:"images"`
	Types   []string               `json:"types"`
	Moves   []string               `json:"moves"`
	Nickame string                 `json:"nickname,omitempty"`
}

type Pokemon struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Types   []string `json:"types"`
	Nickame string   `json:"nickname,omitempty"`
}

type Response struct {
	Message    string
	Code       int
	Percentage int
}

var pokemonDetail []PokemonDetail
var pokemons []Pokemon
var myPokemons []Pokemon
var capturePokemon []Pokemon

func main() {
	r := mux.NewRouter()
	r.Use(commonMiddleware)
	r.HandleFunc("/pokemons", listPokemon).Methods(http.MethodGet)
	r.HandleFunc("/pokemon/{name}", getPokemonDetail).Methods(http.MethodGet)
	r.HandleFunc("/mypokemons", getMyPokemons).Methods(http.MethodGet)
	r.HandleFunc("/capture-pokemon", handleCatch).Methods(http.MethodPost)
	r.HandleFunc("/save-mypokemon", handleMyPokemon).Methods(http.MethodPost)
	r.HandleFunc("/release-pokemon/{id}", releaseMyPokemon).Methods(http.MethodDelete)

	fmt.Println("Server is running on http://localhost:8000")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal("Failed to start server", err)
	}
}

// Fungsi untuk memeriksa apakah sebuah bilangan adalah bilangan prima
func isPrime(number int) bool {
	if number <= 1 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(number))); i++ {
		if number%i == 0 {
			return false
		}
	}
	return true
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func listPokemon(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("https://pokeapi.co/api/v2/pokemon?limit=5")
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer res.Body.Close()
	var data struct {
		Results []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error decoding data:", err)
		return
	}

	for _, pokemon := range data.Results {
		res, err := http.Get(pokemon.URL)
		if err != nil {
			fmt.Println("Error fetching data:", err)
			continue
		}
		url := strings.TrimRight(pokemon.URL, "/")
		_, id := path.Split(url)
		defer res.Body.Close()
		var pokeData struct {
			Sprites struct {
				BackDefault      string `json:"back_default"`
				BackFemale       string `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  string `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"sprites"`
			Types []struct {
				Type struct {
					Name string `json:"name"`
				} `json:"type"`
			} `json:"types"`
			Moves []struct {
				Move struct {
					Name string `json:"name"`
				} `json:"move"`
			} `json:"moves"`
		}
		err = json.NewDecoder(res.Body).Decode(&pokeData)
		if err != nil {
			fmt.Println("Error decoding data:", err)
			continue
		}
		var types []string
		for _, t := range pokeData.Types {
			types = append(types, t.Type.Name)
		}
		var moves []string
		for _, m := range pokeData.Moves {
			moves = append(moves, m.Move.Name)
		}

		idPokemon, _ := strconv.Atoi(id)

		pokemonDetail = append(pokemonDetail, PokemonDetail{
			ID:   idPokemon,
			Name: pokemon.Name,
			Image: map[string]interface{}{
				"back_default":       pokeData.Sprites.BackDefault,
				"back_female":        pokeData.Sprites.BackFemale,
				"back_shiny":         pokeData.Sprites.BackShiny,
				"back_shiny_female":  pokeData.Sprites.BackShiny,
				"front_default":      pokeData.Sprites.FrontDefault,
				"front_female":       pokeData.Sprites.FrontFemale,
				"front_shiny":        pokeData.Sprites.FrontShiny,
				"front_shiny_female": pokeData.Sprites.FrontShinyFemale,
			},
			Types: types,
			Moves: moves,
		})

		pokemons = append(pokemons, Pokemon{
			ID:    idPokemon,
			Name:  pokemon.Name,
			Image: pokeData.Sprites.FrontDefault,
			Types: types,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pokemons)

}

func getPokemonDetail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["name"])
	for _, item := range pokemonDetail {
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Pokemon{})
}

func getMyPokemons(w http.ResponseWriter, r *http.Request) {
	// mengkonversi myPokemonsFromContext menjadi JSON
	myPokemonsJSON, err := json.Marshal(myPokemons)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(myPokemonsJSON)
}

func releaseMyPokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	idPokemon, _ := strconv.Atoi(id)

	n := rand.Intn(21)
	if isPrime(n) {
		for _, pokemon := range myPokemons {
			if pokemon.ID == idPokemon {
				myPokemons = append(myPokemons[:idPokemon-1], myPokemons[idPokemon:]...)
			}
		}
		fmt.Fprintf(w, "Pokemon telah dibebaskan menggunakan bilangan prima %d\n", n)
	} else {
		fmt.Fprintf(w, "Gagal membebaskan Pokemon, bilangan %d bukan bilangan prima\n", n)
	}
}

func handleCatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var catchReq struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&catchReq)
	if err != nil {
		fmt.Println("Error decoding request:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	for i, pokemon := range pokemons {
		if pokemon.Name == catchReq.Name {
			random := rand.Intn(100)
			if random > 50 {
				pokemons[i] = pokemon
				capturePokemon = append(capturePokemon, pokemon)
				w.WriteHeader(http.StatusOK)
				w.Write(jsonResponse("Catch success", http.StatusOK, random))
				return
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(jsonResponse("Catch failed", http.StatusBadRequest, random))
				return
			}
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write(jsonResponse("Pokemon not found", http.StatusBadRequest, 0))
}

func jsonResponse(message string, code int, random int) []byte {
	jsonResponse, err := json.Marshal(Response{Message: message, Code: code, Percentage: random})
	if err != nil {
		fmt.Println("error:", err)
	}
	return jsonResponse
}

func handleMyPokemon(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var addReq struct {
		Nickame string `json:"nickname"`
	}
	err := json.NewDecoder(r.Body).Decode(&addReq)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var newPokemon Pokemon
	for i, pokemon := range capturePokemon {
		newPokemon = Pokemon{
			ID:      i + 1,
			Nickame: addReq.Nickame,
			Name:    pokemon.Name,
			Types:   pokemon.Types,
			Image:   pokemon.Image,
		}
	}

	myPokemons = append(myPokemons, newPokemon)
	json.NewEncoder(w).Encode(newPokemon)
}
