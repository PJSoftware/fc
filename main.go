package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

const apiKeyVar = "API_KEY"

func main() {
	apiKey := readAPIKey()
	fmt.Printf("%s = '%s'\n", apiKeyVar, apiKey)
}

func readAPIKey() string {
	env, err := godotenv.Read()
	if err != nil {
		panic(err)
	}

	if _, ok := env[apiKeyVar]; !ok {
		panic("cannot read " + apiKeyVar + " from .env file")
	}

	apiKey := env[apiKeyVar]
	if apiKey == "xxxx" {
		panic(apiKeyVar + " must be set in .env file")
	}

	return apiKey
}
