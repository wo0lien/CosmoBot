package webhooks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wo0lien/cosmoBot/internal/api"
	"github.com/wo0lien/cosmoBot/internal/logging"
	"github.com/wo0lien/cosmoBot/internal/modules"
	"github.com/wo0lien/cosmoBot/internal/storage/controllers"
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
	// event insert can be a new event or a new link from the dashboard
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

		// check if the event is an insert
		if webhookRes.Type == "records.after.insert" {
			// check if volunteer
			if webhookRes.Data.TableName == "Volunteers" {
				logging.Info.Printf("Received volunteers related webhook.")
				volunteers, err := api.NocoApi.GetAllVolunteers()
				if err != nil {
					logging.Error.Printf("Could not get all volunteers from API. Error: %s", err)
					return
				}
				// check if the volunteer exists in db
				// if not, add it
				vol := controllers.GetVolunteerById(webhookRes.Data.Rows[0].ID)
				if vol == nil {
					controllers.LoadVolunteersToDBFromAPI(volunteers)
				}

				// update volunteer relations
				err = controllers.LoadVolunteersEventsJoinsFromApi(volunteers)
				if err != nil {
					logging.Error.Printf("Could not update volunteers events joins. Error: %s", err)
					return
				}

				// tag all volunteers in all events
				modules.TagAllVolunteersInAllEvents()

				return
			}

			// check if event
			if webhookRes.Data.TableName == "Events" {
				logging.Info.Printf("Received events related webhook.")

				events, err := api.NocoApi.GetAllEvents()
				if err != nil {
					logging.Error.Printf("Could not get all events from API. Error: %s", err)
					return
				}

				volunteers, err := api.NocoApi.GetAllVolunteers()
				if err != nil {
					logging.Error.Printf("Could not get all volunteers from API. Error: %s", err)
					return
				}

				// check if the event exists in db
				// if not, add it
				logging.Debug.Printf("Event id: %d", webhookRes.Data.Rows[0].ID)
				event, err := controllers.EventByID(webhookRes.Data.Rows[0].ID)
				logging.Debug.Printf("Event: %v", event)
				if err != nil {
					logging.Info.Printf("Event does not exist in db. Adding it.\n")
					controllers.LoadEventsInDBFromAPI(*events)
					modules.StartDiscussionForUpcomingEvents()
				}

				// update volunteer relations
				err = controllers.LoadVolunteersEventsJoinsFromApi(volunteers)
				if err != nil {
					logging.Error.Printf("Could not update volunteers events joins. Error: %s", err)
				}

				// tag all volunteers in all events
				modules.TagAllVolunteersInAllEvents()

				return
			}

			logging.Warning.Printf("Received webhook of type %s of unknown table name %s", webhookRes.Type, webhookRes.Data.TableName)
			return
		}
		if webhookRes.Type == "records.after.update" {
			// check if volunteer
			if webhookRes.Data.TableName == "Volunteers" {
				logging.Info.Printf("Received volunteers related webhook.")
				volunteers, err := api.NocoApi.GetAllVolunteers()
				if err != nil {
					logging.Error.Printf("Could not get all volunteers from API. Error: %s", err)
				}
				err = controllers.LoadVolunteersToDBFromAPI(volunteers)

				if err != nil {
					logging.Error.Printf("Could not load volunteers to db. Error: %s", err)
					return
				}

				err = controllers.LoadVolunteersEventsJoinsFromApi(volunteers)

				if err != nil {
					logging.Error.Printf("Could not load volunteers events joins. Error: %s", err)
					return
				}
			}

			// check if event
			if webhookRes.Data.TableName == "Events" {
				logging.Info.Printf("Received events related webhook.")

				events, err := api.NocoApi.GetAllEvents()
				if err != nil {
					logging.Error.Printf("Could not get all events from API. Error: %s", err)
					return
				}

				volunteers, err := api.NocoApi.GetAllVolunteers()
				if err != nil {
					logging.Error.Printf("Could not get all volunteers from API. Error: %s", err)
					return
				}

				err = controllers.LoadEventsInDBFromAPI(*events)
				if err != nil {
					logging.Error.Printf("Could not load events to db. Error: %s", err)
					return
				}

				err = controllers.LoadVolunteersEventsJoinsFromApi(volunteers)

				if err != nil {
					logging.Error.Printf("Could not load volunteers events joins. Error: %s", err)
					return
				}

			}
			return
		}
		logging.Warning.Printf("Received webhook of unknown type %s", webhookRes.Type)

	})

	// Start the HTTP server on port 8080.
	fmt.Println("Server is listening on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
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
