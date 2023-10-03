package api

import "github.com/wo0lien/cosmoBot/internal/logging"

var NocoApi *NocoApiStruct

// Wrapper arround the nocoApi client to used to abstrat the api internal implementation
// TODO remove named field
type NocoApiStruct struct {
	ClientWithResponses *ClientWithResponses
}

func init() {
	NocoApi = NewNocoApi()

	err := NocoApi.ConnectWithEnvApiKey()
	if err != nil {
		logging.Critical.Fatalf("Could not connect to noco api: %s", err)
	}
}

// NewNocoApi create a new NocoApi instance
func NewNocoApi() *NocoApiStruct {
	return &NocoApiStruct{
		ClientWithResponses: nil,
	}
}
