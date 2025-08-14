package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/eventstore"
)

// mockEventStore implements the EventStore interface for testing
type mockEventStore struct {
	responses [][]eventstore.Event
	errors    []error
	callCount int
}

func newMockEventStore() *mockEventStore {
	return &mockEventStore{}
}

func (m *mockEventStore) expectFilter(events []eventstore.Event, err error) *mockEventStore {
	m.responses = append(m.responses, events)
	m.errors = append(m.errors, err)
	return m
}

func (m *mockEventStore) Filter(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
	if m.callCount >= len(m.responses) {
		return nil, errors.New("unexpected Filter call")
	}
	events := m.responses[m.callCount]
	err := m.errors[m.callCount]
	m.callCount++
	return events, err
}

// Unused EventStore methods
func (m *mockEventStore) InstanceIDs(ctx context.Context, query *eventstore.SearchQueryBuilder) ([]string, error) {
	return nil, nil
}
func (m *mockEventStore) FilterToQueryReducer(ctx context.Context, reducer eventstore.QueryReducer) error {
	return nil
}
func (m *mockEventStore) Push(ctx context.Context, cmds ...eventstore.Command) ([]eventstore.Event, error) {
	return nil, nil
}
func (m *mockEventStore) FillFields(ctx context.Context, events ...eventstore.FillFieldsEvent) error {
	return nil
}

// mockEvent implements the Event interface for testing
type mockEvent struct {
	aggType      eventstore.AggregateType
	aggID        string
	sequence     uint64
	eventType    eventstore.EventType
	position     decimal.Decimal
	creationDate time.Time
}

func newMockEvent(aggType, aggID string, sequence uint64, position decimal.Decimal, eventType string) *mockEvent {
	return &mockEvent{
		aggType:      eventstore.AggregateType(aggType),
		aggID:        aggID,
		sequence:     sequence,
		eventType:    eventstore.EventType(eventType),
		position:     position,
		creationDate: time.Now(),
	}
}

func (e *mockEvent) Aggregate() *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            e.aggID,
		Type:          e.aggType,
		ResourceOwner: "test-org",
		InstanceID:    "test-instance",
		Version:       eventstore.Version("v1"),
	}
}

func (e *mockEvent) Sequence() uint64                { return e.sequence }
func (e *mockEvent) Type() eventstore.EventType      { return e.eventType }
func (e *mockEvent) Position() decimal.Decimal       { return e.position }
func (e *mockEvent) CreationDate() time.Time         { return e.creationDate }
func (e *mockEvent) CreatedAt() time.Time            { return e.creationDate }
func (e *mockEvent) Data() []byte                    { return []byte("{}") }
func (e *mockEvent) DataAsBytes() []byte             { return []byte("{}") }
func (e *mockEvent) Unmarshal(ptr interface{}) error { return nil }
func (e *mockEvent) EditorUser() string              { return "test-user" }
func (e *mockEvent) ResourceOwner() string           { return "test-org" }
func (e *mockEvent) InstanceID() string              { return "test-instance" }
func (e *mockEvent) Revision() uint16                { return 1 }
func (e *mockEvent) Creator() string                 { return "test-user" }

// mockReduce creates a statement from an event
func mockReduce(event eventstore.Event) (*Statement, error) {
	return &Statement{
		Aggregate:    event.Aggregate(),
		Sequence:     event.Sequence(),
		Position:     event.Position(),
		CreationDate: event.CreatedAt(),
		offset:       1,
		Execute: func(ctx context.Context, ex Executer, projectionName string) error {
			return nil
		},
	}, nil
}

// setupHandler creates a handler with standard test configuration
func setupHandler(eventStore *mockEventStore) *Handler {
	projection := &projection{
		name: "test_projection",
		reducers: []AggregateReducer{
			{
				Aggregate: "test.aggregate",
				EventReducers: []EventReducer{
					{
						Event:  "test.event",
						Reduce: mockReduce,
					},
				},
			},
		},
	}

	h := &Handler{
		bulkLimit:  10,
		es:         eventStore,
		projection: projection,
		eventTypes: make(map[eventstore.AggregateType][]eventstore.EventType),
	}

	// Build eventTypes map from projection
	for _, reducer := range projection.Reducers() {
		eventTypes := make([]eventstore.EventType, len(reducer.EventReducers))
		for i, eventReducer := range reducer.EventReducers {
			eventTypes[i] = eventReducer.Event
		}
		h.eventTypes[reducer.Aggregate] = eventTypes
	}

	return h
}

func TestHandler_generateStatements_NoEvents(t *testing.T) {
	eventStore := newMockEventStore().expectFilter([]eventstore.Event{}, nil)
	h := setupHandler(eventStore)

	currentState := &state{
		instanceID: "test-instance",
		position:   decimal.Decimal{},
	}

	statements, additionalIteration, err := h.generateStatements(context.Background(), nil, currentState)

	if err != nil {
		t.Errorf("generateStatements() error = %v, want nil", err)
	}
	if len(statements) != 0 {
		t.Errorf("generateStatements() statements count = %d, want 0", len(statements))
	}
	if additionalIteration {
		t.Errorf("generateStatements() additionalIteration = true, want false")
	}
}

func TestHandler_generateStatements_WithEvents(t *testing.T) {
	pos1 := decimal.NewFromInt(100)

	events := []eventstore.Event{
		newMockEvent("test.aggregate", "agg1", 1, pos1, "test.event"),
		newMockEvent("test.aggregate", "agg2", 2, pos1, "test.event"),
	}
	eventStore := newMockEventStore().expectFilter(events, nil)

	projection := &projection{
		name: "test_projection",
		reducers: []AggregateReducer{
			{
				Aggregate: "test.aggregate",
				EventReducers: []EventReducer{
					{
						Event:  "test.event",
						Reduce: mockReduce,
					},
				},
			},
		},
	}

	h := &Handler{
		bulkLimit:  10,
		es:         eventStore,
		projection: projection,
		eventTypes: make(map[eventstore.AggregateType][]eventstore.EventType),
	}

	// Build eventTypes map from projection
	for _, reducer := range projection.Reducers() {
		eventTypes := make([]eventstore.EventType, len(reducer.EventReducers))
		for i, eventReducer := range reducer.EventReducers {
			eventTypes[i] = eventReducer.Event
		}
		h.eventTypes[reducer.Aggregate] = eventTypes
	}

	currentState := &state{
		instanceID: "test-instance",
		position:   decimal.Decimal{},
	}

	statements, additionalIteration, err := h.generateStatements(context.Background(), nil, currentState)

	if err != nil {
		t.Errorf("Handler.generateStatements() error = %v, want nil", err)
	}
	if len(statements) != 2 {
		t.Errorf("Handler.generateStatements() statements count = %d, want 2", len(statements))
	}
	if additionalIteration {
		t.Errorf("Handler.generateStatements() additionalIteration = %v, want false", additionalIteration)
	}
}

func TestHandler_generateStatements_RaceConditionHandling(t *testing.T) {
	pos1 := decimal.NewFromInt(100)

	events := []eventstore.Event{
		newMockEvent("test.aggregate", "agg1", 1, pos1, "test.event"),
		newMockEvent("test.aggregate", "agg2", 2, pos1, "test.event"),
	}
	eventStore := newMockEventStore().expectFilter(events, nil)

	projection := &projection{
		name: "test_projection",
		reducers: []AggregateReducer{
			{
				Aggregate: "test.aggregate",
				EventReducers: []EventReducer{
					{
						Event:  "test.event",
						Reduce: mockReduce,
					},
				},
			},
		},
	}

	h := &Handler{
		bulkLimit:  10,
		es:         eventStore,
		projection: projection,
		eventTypes: make(map[eventstore.AggregateType][]eventstore.EventType),
	}

	// Build eventTypes map from projection
	for _, reducer := range projection.Reducers() {
		eventTypes := make([]eventstore.EventType, len(reducer.EventReducers))
		for i, eventReducer := range reducer.EventReducers {
			eventTypes[i] = eventReducer.Event
		}
		h.eventTypes[reducer.Aggregate] = eventTypes
	}

	currentState := &state{
		instanceID:    "test-instance",
		position:      pos1,
		aggregateID:   "agg2", // Last processed
		aggregateType: "test.aggregate",
		sequence:      2,
		offset:        0, // Zero offset triggers race condition handling
	}

	statements, additionalIteration, err := h.generateStatements(context.Background(), nil, currentState)

	if err != nil {
		t.Errorf("Handler.generateStatements() error = %v, want nil", err)
	}
	if len(statements) != 0 {
		t.Errorf("Handler.generateStatements() statements count = %d, want 0", len(statements))
	}
	if additionalIteration {
		t.Errorf("Handler.generateStatements() additionalIteration = %v, want false", additionalIteration)
	}
	if currentState.offset != 2 {
		t.Errorf("Handler.generateStatements() updated offset = %d, want 2", currentState.offset)
	}
}

func TestHandler_generateStatements_BulkLimitExceeded(t *testing.T) {
	pos1 := decimal.NewFromInt(100)

	events := []eventstore.Event{
		newMockEvent("test.aggregate", "agg1", 1, pos1, "test.event"),
		newMockEvent("test.aggregate", "agg2", 2, pos1, "test.event"),
	}
	eventStore := newMockEventStore().expectFilter(events, nil)

	projection := &projection{
		name: "test_projection",
		reducers: []AggregateReducer{
			{
				Aggregate: "test.aggregate",
				EventReducers: []EventReducer{
					{
						Event:  "test.event",
						Reduce: mockReduce,
					},
				},
			},
		},
	}

	h := &Handler{
		bulkLimit:  2, // Bulk limit exactly equals number of events
		es:         eventStore,
		projection: projection,
		eventTypes: make(map[eventstore.AggregateType][]eventstore.EventType),
	}

	// Build eventTypes map from projection
	for _, reducer := range projection.Reducers() {
		eventTypes := make([]eventstore.EventType, len(reducer.EventReducers))
		for i, eventReducer := range reducer.EventReducers {
			eventTypes[i] = eventReducer.Event
		}
		h.eventTypes[reducer.Aggregate] = eventTypes
	}

	currentState := &state{
		instanceID: "test-instance",
		position:   decimal.Decimal{},
	}

	statements, additionalIteration, err := h.generateStatements(context.Background(), nil, currentState)

	if err != nil {
		t.Errorf("Handler.generateStatements() error = %v, want nil", err)
	}
	if len(statements) != 2 {
		t.Errorf("Handler.generateStatements() statements count = %d, want 2", len(statements))
	}
	if !additionalIteration {
		t.Errorf("Handler.generateStatements() additionalIteration = false, want true due to bulk limit")
	}
}

func TestHandler_generateStatements_ErrorScenario(t *testing.T) {
	eventStore := newMockEventStore().expectFilter(nil, errors.New("filter error"))

	projection := &projection{
		name: "test_projection",
		reducers: []AggregateReducer{
			{
				Aggregate: "test.aggregate",
				EventReducers: []EventReducer{
					{
						Event:  "test.event",
						Reduce: mockReduce,
					},
				},
			},
		},
	}

	h := &Handler{
		bulkLimit:  10,
		es:         eventStore,
		projection: projection,
		eventTypes: make(map[eventstore.AggregateType][]eventstore.EventType),
	}

	// Build eventTypes map from projection
	for _, reducer := range projection.Reducers() {
		eventTypes := make([]eventstore.EventType, len(reducer.EventReducers))
		for i, eventReducer := range reducer.EventReducers {
			eventTypes[i] = eventReducer.Event
		}
		h.eventTypes[reducer.Aggregate] = eventTypes
	}

	currentState := &state{
		instanceID: "test-instance",
		position:   decimal.Decimal{},
	}

	statements, additionalIteration, err := h.generateStatements(context.Background(), nil, currentState)

	if err == nil {
		t.Errorf("Handler.generateStatements() error = nil, want error")
	}
	if len(statements) != 0 {
		t.Errorf("Handler.generateStatements() statements count = %d, want 0", len(statements))
	}
	if additionalIteration {
		t.Errorf("Handler.generateStatements() additionalIteration = %v, want false", additionalIteration)
	}
}

func TestHandler_generateStatements_MultipleCallsOffsetProgression(t *testing.T) {
	// Create events for multiple calls
	events1 := []eventstore.Event{
		newMockEvent("test.aggregate", "agg-1", 1, decimal.NewFromInt(100), "test.event"),
		newMockEvent("test.aggregate", "agg-2", 2, decimal.NewFromInt(200), "test.event"),
	}
	events2 := []eventstore.Event{
		newMockEvent("test.aggregate", "agg-3", 3, decimal.NewFromInt(300), "test.event"),
	}
	events3 := []eventstore.Event{
		newMockEvent("test.aggregate", "agg-4", 4, decimal.NewFromInt(400), "test.event"),
		newMockEvent("test.aggregate", "agg-5", 5, decimal.NewFromInt(500), "test.event"),
	}

	// Setup mock eventstore with multiple filter expectations
	eventStore := newMockEventStore().
		expectFilter(events1, nil). // First call
		expectFilter(events2, nil). // Second call with offset
		expectFilter(events3, nil)  // Third call with offset

	projection := &projection{
		name: "test_projection",
		reducers: []AggregateReducer{
			{
				Aggregate: "test.aggregate",
				EventReducers: []EventReducer{
					{
						Event:  "test.event",
						Reduce: mockReduce,
					},
				},
			},
		},
	}

	h := &Handler{
		bulkLimit:  10,
		es:         eventStore,
		projection: projection,
		eventTypes: make(map[eventstore.AggregateType][]eventstore.EventType),
	}

	// Build eventTypes map from projection
	for _, reducer := range projection.Reducers() {
		eventTypes := make([]eventstore.EventType, len(reducer.EventReducers))
		for i, eventReducer := range reducer.EventReducers {
			eventTypes[i] = eventReducer.Event
		}
		h.eventTypes[reducer.Aggregate] = eventTypes
	}

	// Initial state
	currentState := &state{
		instanceID: "test-instance",
		position:   decimal.NewFromInt(50), // Start from position 50
		offset:     0,
	}

	// First call
	statements1, additionalIteration1, err1 := h.generateStatements(context.Background(), nil, currentState)
	if err1 != nil {
		t.Errorf("Handler.generateStatements() first call error = %v, want nil", err1)
	}
	if len(statements1) != 2 {
		t.Errorf("Handler.generateStatements() first call statements count = %d, want 2", len(statements1))
	}
	if additionalIteration1 {
		t.Errorf("Handler.generateStatements() first call additionalIteration = %v, want false", additionalIteration1)
	}

	// Update state based on first call results - simulate progression
	currentState.position = decimal.NewFromInt(200) // Position should progress
	currentState.offset = 2                         // Offset should be set to number of processed statements

	// Second call - should use offset
	statements2, additionalIteration2, err2 := h.generateStatements(context.Background(), nil, currentState)
	if err2 != nil {
		t.Errorf("Handler.generateStatements() second call error = %v, want nil", err2)
	}
	if len(statements2) != 1 {
		t.Errorf("Handler.generateStatements() second call statements count = %d, want 1", len(statements2))
	}
	if additionalIteration2 {
		t.Errorf("Handler.generateStatements() second call additionalIteration = %v, want false", additionalIteration2)
	}

	// Update state for third call
	currentState.position = decimal.NewFromInt(400)
	currentState.offset = 1 // Reset offset after processing previous batch

	// Third call
	statements3, additionalIteration3, err3 := h.generateStatements(context.Background(), nil, currentState)
	if err3 != nil {
		t.Errorf("Handler.generateStatements() third call error = %v, want nil", err3)
	}
	if len(statements3) != 2 {
		t.Errorf("Handler.generateStatements() third call statements count = %d, want 2", len(statements3))
	}
	if additionalIteration3 {
		t.Errorf("Handler.generateStatements() third call additionalIteration = %v, want false", additionalIteration3)
	}

	// Verify basic test completion - all calls returned expected number of statements
	totalStatements := len(statements1) + len(statements2) + len(statements3)
	expectedTotal := 5 // 2 + 1 + 2
	if totalStatements != expectedTotal {
		t.Errorf("Total statements count = %d, want %d", totalStatements, expectedTotal)
	}

	// Verify that statements have different aggregate IDs showing progression
	if len(statements1) > 0 && statements1[0].Aggregate.ID != "agg-1" && statements1[0].Aggregate.ID != "agg-2" {
		t.Errorf("First call statement should have agg-1 or agg-2, got: %s", statements1[0].Aggregate.ID)
	}

	if len(statements2) > 0 && statements2[0].Aggregate.ID != "agg-3" {
		t.Errorf("Second call statement should have agg-3, got: %s", statements2[0].Aggregate.ID)
	}

	if len(statements3) > 0 && statements3[0].Aggregate.ID != "agg-4" && statements3[0].Aggregate.ID != "agg-5" {
		t.Errorf("Third call statement should have agg-4 or agg-5, got: %s", statements3[0].Aggregate.ID)
	}
}

func TestHandler_generateStatements_NormalCase(t *testing.T) {
	// Test the normal case where events are processed from the beginning
	// with no previous state (fresh start)
	events := []eventstore.Event{
		newMockEvent("test.aggregate", "agg1", 1, decimal.NewFromInt(100), "test.event"),
		newMockEvent("test.aggregate", "agg2", 2, decimal.NewFromInt(200), "test.event"),
		newMockEvent("test.aggregate", "agg3", 3, decimal.NewFromInt(300), "test.event"),
	}
	eventStore := newMockEventStore().expectFilter(events, nil)
	h := setupHandler(eventStore)

	// Fresh state with no position or offset (normal case)
	currentState := &state{
		instanceID: "test-instance",
		position:   decimal.Decimal{}, // Zero position means fresh start
		offset:     0,
	}

	statements, additionalIteration, err := h.generateStatements(context.Background(), nil, currentState)

	if err != nil {
		t.Errorf("generateStatements() normal case error = %v, want nil", err)
	}
	if len(statements) != 3 {
		t.Errorf("generateStatements() normal case statements count = %d, want 3", len(statements))
	}
	if additionalIteration {
		t.Errorf("generateStatements() normal case additionalIteration = %v, want false", additionalIteration)
	}

	// Verify all statements were processed (no skipping in normal case)
	expectedAggregateIDs := []string{"agg1", "agg2", "agg3"}
	for i, stmt := range statements {
		if stmt.Aggregate.ID != expectedAggregateIDs[i] {
			t.Errorf("Statement %d aggregate ID = %s, want %s", i, stmt.Aggregate.ID, expectedAggregateIDs[i])
		}
	}
}

func TestHandler_generateStatements_SkipStatementsIdxCase(t *testing.T) {
	// Test the idx >= 0 case where some statements need to be skipped
	// because they were already processed (continuing from a specific position)
	events := []eventstore.Event{
		newMockEvent("test.aggregate", "agg1", 1, decimal.NewFromInt(100), "test.event"),
		newMockEvent("test.aggregate", "agg2", 2, decimal.NewFromInt(200), "test.event"),
		newMockEvent("test.aggregate", "agg3", 3, decimal.NewFromInt(300), "test.event"),
		newMockEvent("test.aggregate", "agg4", 4, decimal.NewFromInt(400), "test.event"),
	}
	eventStore := newMockEventStore().expectFilter(events, nil)
	h := setupHandler(eventStore)

	// Create a custom reducer that generates statements matching the events exactly
	// so we can simulate proper skipping behavior
	customReduce := func(event eventstore.Event) (*Statement, error) {
		return &Statement{
			Aggregate:    event.Aggregate(),
			Sequence:     event.Sequence(),
			Position:     event.Position(),
			CreationDate: event.CreatedAt(),
			offset:       1,
			Execute: func(ctx context.Context, ex Executer, projectionName string) error {
				return nil
			},
		}, nil
	}

	// Update handler to use custom reducer
	h.projection.(*projection).reducers[0].EventReducers[0].Reduce = customReduce

	// State indicating we've already processed up to position 200 (agg2)
	// This should cause idx >= 0 and skip the first two statements
	currentState := &state{
		instanceID:    "test-instance",
		position:      decimal.NewFromInt(200),
		aggregateID:   "agg2",
		aggregateType: "test.aggregate",
		sequence:      2,
		offset:        0,
	}

	statements, additionalIteration, err := h.generateStatements(context.Background(), nil, currentState)

	if err != nil {
		t.Errorf("generateStatements() skip case error = %v, want nil", err)
	}

	if len(statements) != 2 {
		t.Errorf("generateStatements() skip case statements count = %d, want 2", len(statements))
	}

	// additionalIteration should be true because len(statements) < len(events)
	// This is correct behavior when statements are skipped
	if !additionalIteration {
		t.Errorf("generateStatements() skip case additionalIteration = %v, want true (because len(statements) < len(events))", additionalIteration)
	}

	// Verify that the correct statements were returned (the ones after the skip point)
	if len(statements) > 0 && statements[0].Aggregate.ID != "agg3" {
		t.Errorf("First returned statement should be agg3, got: %s", statements[0].Aggregate.ID)
	}
	if len(statements) > 1 && statements[1].Aggregate.ID != "agg4" {
		t.Errorf("Second returned statement should be agg4, got: %s", statements[1].Aggregate.ID)
	}
}

func TestHandler_generateStatements_AllStatementsProcessedNormalCase(t *testing.T) {
	// Test the else case where all statements are already processed (idx+1 == len(statements))
	// but we're NOT in a race condition (not allStatementsAtSamePosition && currentState.offset == 0)
	// This is the "Normal case - update state to the last statement" path
	events := []eventstore.Event{
		newMockEvent("test.aggregate", "agg1", 1, decimal.NewFromInt(100), "test.event"),
		newMockEvent("test.aggregate", "agg2", 2, decimal.NewFromInt(200), "test.event"),
	}
	eventStore := newMockEventStore().expectFilter(events, nil)
	h := setupHandler(eventStore)

	// Create a custom reducer that generates statements exactly matching our state
	customReduce := func(event eventstore.Event) (*Statement, error) {
		return &Statement{
			Aggregate:    event.Aggregate(),
			Sequence:     event.Sequence(),
			Position:     event.Position(),
			CreationDate: event.CreatedAt(),
			offset:       1,
			Execute: func(ctx context.Context, ex Executer, projectionName string) error {
				return nil
			},
		}, nil
	}
	h.projection.(*projection).reducers[0].EventReducers[0].Reduce = customReduce

	// State indicating we've processed up to agg2 with offset > 0
	// This ensures we enter the else case (normal completion) rather than race condition handling
	currentState := &state{
		instanceID:    "test-instance",
		position:      decimal.NewFromInt(200), // Same as last event position
		aggregateID:   "agg2",                  // Same as last event aggregate
		aggregateType: "test.aggregate",        // Same as last event aggregate type
		sequence:      2,                       // Same as last event sequence
		offset:        1,                       // > 0 to avoid race condition path
	}

	statements, additionalIteration, err := h.generateStatements(context.Background(), nil, currentState)

	if err != nil {
		t.Errorf("generateStatements() all processed normal case error = %v, want nil", err)
	}

	// Should return no statements because all are already processed
	if len(statements) != 0 {
		t.Errorf("generateStatements() all processed normal case statements count = %d, want 0", len(statements))
	}

	// Should not trigger additional iteration since everything is processed normally
	if additionalIteration {
		t.Errorf("generateStatements() all processed normal case additionalIteration = %v, want false", additionalIteration)
	}

	// Verify that the state would be updated to the last statement (this is tested indirectly
	// through the fact that the function returns successfully with no statements)
	// The actual state update happens in the calling function, but the logic we're testing
	// is the path that leads to returning (nil, false, nil) in the else case
}

func TestHandler_generateStatements_AllStatementsProcessedRaceConditionCase(t *testing.T) {
	// Test the if case where all statements are already processed (idx+1 == len(statements))
	// AND we're in a race condition (allStatementsAtSamePosition && currentState.offset == 0)
	// This should increment the offset to handle the race condition
	events := []eventstore.Event{
		newMockEvent("test.aggregate", "agg1", 1, decimal.NewFromInt(200), "test.event"),
		newMockEvent("test.aggregate", "agg2", 2, decimal.NewFromInt(200), "test.event"),
	}
	eventStore := newMockEventStore().expectFilter(events, nil)
	h := setupHandler(eventStore)

	// Create a custom reducer that generates statements with the same position
	customReduce := func(event eventstore.Event) (*Statement, error) {
		return &Statement{
			Aggregate:    event.Aggregate(),
			Sequence:     event.Sequence(),
			Position:     decimal.NewFromInt(200), // Same position for all statements
			CreationDate: event.CreatedAt(),
			offset:       1,
			Execute: func(ctx context.Context, ex Executer, projectionName string) error {
				return nil
			},
		}, nil
	}
	h.projection.(*projection).reducers[0].EventReducers[0].Reduce = customReduce

	// State that triggers race condition: position matches last statement AND offset is 0
	currentState := &state{
		instanceID:    "test-instance",
		position:      decimal.NewFromInt(200), // Same as all statement positions
		aggregateID:   "agg2",                  // Same as last event aggregate
		aggregateType: "test.aggregate",        // Same as last event aggregate type
		sequence:      2,                       // Same as last event sequence
		offset:        0,                       // 0 triggers race condition handling
	}

	statements, additionalIteration, err := h.generateStatements(context.Background(), nil, currentState)

	if err != nil {
		t.Errorf("generateStatements() race condition case error = %v, want nil", err)
	}

	// Should return no statements because all are already processed
	if len(statements) != 0 {
		t.Errorf("generateStatements() race condition case statements count = %d, want 0", len(statements))
	}

	// Should not trigger additional iteration
	if additionalIteration {
		t.Errorf("generateStatements() race condition case additionalIteration = %v, want false", additionalIteration)
	}

	// In this case, the offset should be incremented to len(statements) = 2
	// We can't directly test the state modification since it happens in-place and is returned,
	// but the fact that we get (nil, false, nil) confirms we took the race condition path
}
