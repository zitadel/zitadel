package requesttracker

import (
	"context"
	"net/http"

	"github.com/crewjam/saml/samlsp"
)

type GetRequest func(ctx context.Context, intentID string) (*samlsp.TrackedRequest, error)
type AddRequest func(ctx context.Context, intentID, requestID string) error

type RequestTracker struct {
	addRequest AddRequest
	getRequest GetRequest
}

func New(addRequestF AddRequest, getRequestF GetRequest) samlsp.RequestTracker {
	return &RequestTracker{
		addRequest: addRequestF,
		getRequest: getRequestF,
	}
}

func (rt *RequestTracker) TrackRequest(w http.ResponseWriter, r *http.Request, samlRequestID string) (index string, err error) {
	// intentID is stored in r.URL
	intentID := r.URL.String()
	if err := rt.addRequest(r.Context(), intentID, samlRequestID); err != nil {
		return "", err
	}
	return intentID, nil
}

func (rt *RequestTracker) StopTrackingRequest(w http.ResponseWriter, r *http.Request, index string) error {
	// error is not handled in SP logic
	return nil
}

func (rt *RequestTracker) GetTrackedRequests(r *http.Request) []samlsp.TrackedRequest {
	// RelayState is the context of the auth flow and as such contains the intentID
	intentID := r.FormValue("RelayState")

	request, err := rt.getRequest(r.Context(), intentID)
	if err != nil {
		return nil
	}
	return []samlsp.TrackedRequest{
		{
			Index:         request.Index,
			SAMLRequestID: request.SAMLRequestID,
			URI:           request.URI,
		},
	}
}

func (rt *RequestTracker) GetTrackedRequest(r *http.Request, index string) (*samlsp.TrackedRequest, error) {
	return rt.getRequest(r.Context(), index)
}
