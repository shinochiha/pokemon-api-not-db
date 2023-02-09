package models

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
)

type Pokemon struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Nickname string  `json:"nickname"`
	Height   float32 `json:"height"`
	Weight   float32 `json:"weight"`
	Sprites  struct {
		BackDefault  string `json:"back_default"`
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
}

type MyPokemon struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Nickname string  `json:"nickname"`
	Height   float32 `json:"height"`
	Weight   float32 `json:"weight"`
	Sprites  struct {
		BackDefault  string `json:"back_default"`
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
}

type PokemonDetail struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Height  float32 `json:"height"`
	Weight  float32 `json:"weight"`
	Sprites struct {
		BackDefault  string `json:"back_default"`
		FrontDefault string `json:"front_default"`
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

var pokemonsDetail []PokemonDetail
var pokemonURL PokemonDetail
var myPokemons []MyPokemon

// PokemonList represents a list of pokemon
type PokemonList struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []Pokemon `json:"results"`
}

func GetPokemonList() (PokemonList, error) {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon?limit=100")
	if err != nil {
		return PokemonList{}, err
	}
	defer resp.Body.Close()

	var data struct {
		Results []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return PokemonList{}, err
	}

	var pokemons []Pokemon
	for i, result := range data.Results {
		resp, err := http.Get(result.URL)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		// var pokemonURL Pokemon
		err = json.Unmarshal(body, &pokemonURL)
		if err != nil {
			panic(err)
		}

		pokemonsDetail = append(pokemonsDetail, PokemonDetail{
			Name:    result.Name,
			ID:      i + 1,
			Height:  pokemonURL.Height,
			Weight:  pokemonURL.Weight,
			Sprites: pokemonURL.Sprites,
			Types:   pokemonURL.Types,
			Moves:   pokemonURL.Moves,
		})

		pokemon := Pokemon{
			Name:    result.Name,
			ID:      i + 1,
			Height:  pokemonURL.Height,
			Weight:  pokemonURL.Weight,
			Sprites: pokemonURL.Sprites,
		}
		pokemons = append(pokemons, pokemon)
	}

	return PokemonList{Results: pokemons}, nil
}

func GetPokemonByID(id int) (PokemonDetail, error) {
	resp, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", id))
	if err != nil {
		return PokemonDetail{}, err
	}
	defer resp.Body.Close()

	var data PokemonDetail
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return PokemonDetail{}, err
	}

	return PokemonDetail{ID: data.ID, Name: data.Name, Height: data.Height, Weight: data.Weight, Sprites: data.Sprites, Types: data.Types, Moves: data.Moves}, nil
}

func CatchPokemonById(pokemonID int) (MyPokemon, error) {
	successRate := rand.Intn(100)
	if successRate < 50 {
		// kalau gagal menangkap pokemon
		return MyPokemon{}, fmt.Errorf("failed to catch pokemon with id %d because success rate is %d", pokemonID, successRate)
	}

	// kalau berhasil menangkap pokemon
	pokemon, err := GetPokemonByID(pokemonID)
	if err != nil {
		return MyPokemon{}, err
	}

	fmt.Print("Give a nickname for your new pokemon: ")
	var nickname string
	_, err = fmt.Scan(&nickname)
	if err != nil {
		return MyPokemon{}, err
	}

	myPokemon := MyPokemon{
		ID:       pokemon.ID,
		Name:     pokemon.Name,
		Nickname: nickname,
		Height:   pokemon.Height,
		Weight:   pokemon.Weight,
		Sprites:  pokemon.Sprites,
	}

	myPokemons = append(myPokemons, myPokemon)
	return myPokemon, nil
}

// func GetMyPokemon() []MyPokemon {
// 	return myPokemons
// }

func GetMyPokemon() (int, []MyPokemon) {
	return len(myPokemons), myPokemons
}

func ReleasePokemon(id int, primeNumber int) bool {
	released := false
	limit := int(math.Sqrt(float64(primeNumber))) + 1
	for i := 2; i < limit; i++ {
		if primeNumber%i == 0 {
			released = true
			break
		}
	}

	if primeNumber == 1 {
		released = true
	}

	if !released {
		for i, pokemon := range myPokemons {
			if pokemon.ID == id {
				myPokemons = append(myPokemons[:i], myPokemons[i+1:]...)
				return true
			}
		}
		return false
	}
	return false
}

func ChangeNickname(id int, nickname string) error {
	for i, pokemon := range myPokemons {
		if pokemon.ID == id {
			fib := fibonacci()
			myPokemons[i].Nickname = fmt.Sprintf("%s-%d", nickname, fib)
			return nil
		}

	}
	return fmt.Errorf("Pokemon with ID %d not found", id)
}

var count int

func fibonacci() int {
	count++
	if count == 1 {
		return 0
	}
	if count == 2 {
		return 1
	}
	if count == 3 {
		return 1
	}
	if count == 4 {
		return 2
	}
	if count == 5 {
		return 3
	}
	if count == 6 {
		return 5
	}
	if count == 7 {
		return 8
	}
	return fibonacci() + fibonacci()
}
