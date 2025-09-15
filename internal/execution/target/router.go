package target

import (
	"slices"
	"strings"
)

type element struct {
	ID      string   `json:"id"`
	Targets []Target `json:"targets,omitempty"`
}

type Router []element

func NewRouter(targets []Target) Router {
	m := make(map[string][]Target)
	for _, t := range targets {
		m[t.GetExecutionID()] = append(m[t.GetExecutionID()], t)
	}
	router := make(Router, 0, len(m))
	for id, targets := range m {
		router = append(router, element{
			ID:      id,
			Targets: targets,
		})
	}
	slices.SortFunc(router, func(a, b element) int {
		return strings.Compare(a.ID, b.ID)
	})
	return router
}

// Get execution targets by exact match of the executionID
func (r Router) Get(executionID string) ([]Target, bool) {
	i, ok := slices.BinarySearchFunc(r, executionID, func(a element, b string) int {
		return strings.Compare(a.ID, b)
	})
	if ok {
		return r[i].Targets, true
	}
	return nil, false
}

// GetEventBestMatch returns the best matching execution targets for an event.
// The following match priority is used:
//  1. Exact match
//  2. Wildcard match
//  3. Prefix match ("event")
func (r Router) GetEventBestMatch(executionID string) ([]Target, bool) {
	t, ok := r.Get(executionID)
	if ok {
		return t, true
	}
	var bestMatch element
	for _, e := range r {
		if e.ID == "event" && strings.HasPrefix(executionID, e.ID) {
			bestMatch, ok = e, true
		}
		cut, has := strings.CutSuffix(e.ID, ".*")
		if has && strings.HasPrefix(executionID, cut) {
			bestMatch, ok = e, true
		}
	}
	return bestMatch.Targets, ok
}

func (r Router) IsZero() bool {
	return len(r) == 0
}
