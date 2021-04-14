package eventstore

import "time"

type Handler struct {
	Eventstore   *Eventstore
	Queue        *JobQueue
	RequeueAfter time.Duration
	Timer        *time.Timer
	Sub          *Subscription
}

func NewHandler(
	eventstore *Eventstore,
	queue *JobQueue,
	requeueAfter time.Duration,
	aggregates ...AggregateType,
) Handler {
	h := Handler{
		Eventstore:   eventstore,
		Queue:        queue,
		RequeueAfter: requeueAfter,
		// first requeue is instant on startup
		Timer: time.NewTimer(0),
	}

	return h
}

func (h Handler) Subscribe(aggregates ...AggregateType) {
	h.Sub = Subscribe(aggregates...)
}

func (h Handler) ResetTimer() {
	h.Timer.Reset(h.RequeueAfter)
}
