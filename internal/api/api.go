package api

import (
	"io"
	"net/http"
	"os"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/joho/godotenv"
	"github.com/wo0lien/cosmoBot/internal/logging"
)

// ########## CONSTANTS ##########

var NOCO_API_KEY string
var NOCO_URL string

// ########## STRUCTS ##########

var NocoApi *NocoApiStruct

// Wrapper arround the nocoApi client to used to abstrat the api internal implementation
// TODO remove named field
type NocoApiStruct struct {
	clientWithResponses *ClientWithResponses
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
		clientWithResponses: nil,
	}
}

// date in nocodb are in this format
// use this layout in time.Parse function to get go time.Time matching object
var NOCO_TIME_LAYOUT = "2006-01-02 15:04:05-07:00"

// ConnectWithEnvApiKey connect to the noco api using the api key from the environment variable NOCO_API_KEY
func (na *NocoApiStruct) ConnectWithEnvApiKey() error {
	// load the api key from the environment variable
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	// load env variables BOT_TOKEN
	NOCO_API_KEY = os.Getenv("NOCO_API_KEY")
	if NOCO_API_KEY == "" {
		panic("NOCO_API_KEY env variable is not set")
	}

	NOCO_URL = os.Getenv("NOCO_URL")
	if NOCO_URL == "" {
		panic("NOCO_URL env variable is not set")
	}

	return na.Connect(NOCO_URL, NOCO_API_KEY)
}

// Connect to the noco api using the given api key
func (na *NocoApiStruct) Connect(url, apiKey string) error {

	apiKeyProvider, err := securityprovider.NewSecurityProviderApiKey("header", "xc-token", apiKey)

	if err != nil {
		return err
	}

	client, err := NewClientWithResponses(url, WithRequestEditorFn(apiKeyProvider.Intercept))

	if err != nil {
		return err
	}

	na.clientWithResponses = client

	return nil
}

// GetBodyString return the body of the given http response as a string
func GetBodyString(res http.Response) (string, error) {
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}
