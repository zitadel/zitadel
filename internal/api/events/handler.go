package events

import (
	"encoding/json"
	"net/http"

	"github.com/rs/cors"
	"github.com/zitadel/zitadel/internal/database"
)

type eventRequest struct {
	InstanceID string `json:"instance_id"`
	ParentType string `json:"parent_type"`
	ParentID   string `json:"parent_id"`
	TableName  string `json:"table_name"`
	Event      string `json:"event"`
}

// Start returns an http.Handler that accepts POST / to insert a row into
// projections.service_ping_resource_events for quick debugging/telemetry.
func Start(db *database.DB) (http.Handler, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()
		var payload eventRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}
		if payload.InstanceID == "" || payload.ParentType == "" || payload.ParentID == "" {
			http.Error(w, "instance_id, parent_type and parent_id are required", http.StatusBadRequest)
			return
		}
		if payload.TableName == "" {
			payload.TableName = "custom.events"
		}
		if payload.Event == "" {
			payload.Event = "CUSTOM"
		}
		_, err := db.ExecContext(r.Context(),
			`INSERT INTO projections.service_ping_resource_events (instance_id, table_name, parent_type, parent_id, event)
			 VALUES ($1, $2, $3, $4, $5)`,
			payload.InstanceID, payload.TableName, payload.ParentType, payload.ParentID, payload.Event,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	return cors.AllowAll().Handler(mux), nil
}
