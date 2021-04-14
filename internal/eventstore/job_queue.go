package eventstore

import (
	"context"
	"os"
	"strconv"
	"sync"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/id"
)

type JobQueue struct {
	queue      chan *job
	cancel     context.CancelFunc
	es         *Eventstore
	shutdownWG sync.WaitGroup
}

//NewJobQueue creates a new job queue
// and sets up the workers which will process the queued jobs
func NewJobQueue(ctx context.Context, workerCount uint16, eventstore *Eventstore) *JobQueue {
	lockerID, err := os.Hostname()
	if err != nil || lockerID == "" {
		lockerID, err = id.SonyFlakeGenerator.Next()
		logging.Log("SPOOL-bdO56").OnError(err).Panic("unable to generate lockID")
	}

	defer logging.LogWithFields("EVENT-4DMkn", "lockerID", lockerID, "workerCount", workerCount).Info("job queue startedstarted")

	ctx, cancel := context.WithCancel(ctx)
	queue := &JobQueue{
		es:     eventstore,
		cancel: cancel,
		//TODO: optimise size of queue
		queue:      make(chan *job, 50),
		shutdownWG: sync.WaitGroup{},
	}

	for i := uint16(0); i < workerCount; i++ {
		workerID := lockerID + "-" + strconv.Itoa(int(i))
		queue.shutdownWG.Add(1)
		go queue.startWorker(ctx, workerID)
	}
	return queue
}

//Queue adds a job to the internal channel
// if the latest sequence of the passed query is higher than existing.
// The caller has to validate if the returned sequence is higher that the requested one
// if no, no events will be emitted to his queue
func (q *JobQueue) Queue(ctx context.Context, query *SearchQueryBuilder, queue chan EventReader) (latestSequence uint64, err error) {
	seq, err := q.es.LatestSequence(ctx, query)
	if err != nil {
		return 0, err
	}
	if query.eventSequence < seq {
		q.queue <- &job{ctx: ctx, query: query, eventQueue: queue}
	}
	return seq, nil
}

//Shutdown gracefully stops all workers
// and returns as soon as all workers are stopped
func (q *JobQueue) Shutdown() {
	q.cancel()
	q.shutdownWG.Wait()
}

func (q *JobQueue) execute(j *job) {
	events, err := q.es.FilterEvents(j.ctx, j.query)
	if err != nil {
		logging.Log("EVENT-BcVzZ").WithError(err).Warn("unable to load events")
		return
	}
	for _, event := range events {
		j.eventQueue <- event
	}
}

func (q *JobQueue) startWorker(ctx context.Context, workerID string) {
	defer logging.LogWithFields("EVENT-LAIyN", "workerID", workerID).Warn("Worker shut down")
	for {
		select {
		case <-ctx.Done():
			q.shutdownWG.Done()
			return
		case job := <-q.queue:
			q.execute(job)
		}
	}
}

type job struct {
	ctx        context.Context
	query      *SearchQueryBuilder
	eventQueue chan EventReader
}
