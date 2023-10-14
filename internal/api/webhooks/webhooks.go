package webhooks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wo0lien/cosmoBot/internal/calendar"
	"github.com/wo0lien/cosmoBot/internal/logging"
	"github.com/wo0lien/cosmoBot/internal/workflows"
)

type WebHookResponseDataRow struct {
	ID uint `json:"id"`
	// + other fields that are not relevant for us
}

type WebHookResponseData struct {
	TableName string                   `json:"table_name"`
	Rows      []WebHookResponseDataRow `json:"rows"`
	// + other fields that are not relevant for us
}

type WebHookResponse struct {
	Data WebHookResponseData `json:"data"`
	Type string              `json:"type"`
}

// Handle websockets

func StartWebHooksHandlingServer() {
	// Get the current calendar service
	cs, err := calendar.CalendarService()
	if err != nil {
		panic(err)
	}
	// event insert can be a new event or a new link from the dashboard
	http.HandleFunc("/", httpHandler(cs))

	// Start the HTTP server on port 8080.
	fmt.Println("Server is listening on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

// wrapper around the http handler to pass the calendar service
func httpHandler(cs *calendar.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Print the request method and path to the console.
		logging.Info.Printf("Webhook server received %s request for %s\n", r.Method, r.URL.Path)

		// Write a response back to the client.
		w.WriteHeader(http.StatusOK)

		webhookRes, err := parseBody(w, r)

		// do not continue in case of error
		if err != nil {
			logging.Warning.Printf("Could not parse webhook data. Stop processing this webhook call. Error : %s", err)
			return
		}

		identifierTypeTableName := webhookRes.Type + webhookRes.Data.TableName
		rowId := webhookRes.Data.Rows[0].ID

		logging.Info.Printf("Handling webhook of type %s", identifierTypeTableName)

		switch identifierTypeTableName {
		case "records.after.insertVolunteers":
			err = workflows.InsertVolunteerByID(cs, rowId)
		case "records.after.updateVolunteers":
			err = workflows.UpdateVolunteerByID(cs, rowId)
		case "records.after.deleteVolunteers":
			err = workflows.DeleteVolunteerByID(rowId)
		case "records.after.insertEvents":
			err = workflows.InsertEventByID(cs, rowId)
		case "records.after.updateEvents":
			err = workflows.UpdateEventByID(cs, rowId)
		case "records.after.deleteEvents":
			err = workflows.DeleteEventByID(cs, rowId)
		default:
			logging.Warning.Printf("Received webhook of unknown type %s", identifierTypeTableName)
			return
		}
		// logging error
		if err != nil {
			logging.Warning.Printf("Could not handle webhook of type %s. Error: %s", identifierTypeTableName, err)
		}
	}
}

func parseBody(w http.ResponseWriter, r *http.Request) (*WebHookResponse, error) {
	// print body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var data WebHookResponse

	err = json.Unmarshal(body, &data)
	if err != nil {

		return nil, err
	}

	return &data, nil
}
