package queue

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// JobQueue represents a job queue for processing tasks asynchronously
type JobQueue struct {
	jobs    chan Job
	workers int
	logger  *slog.Logger
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

// Job represents a single job to be processed
type Job struct {
	ID       string
	Task     func() error
	RetryMax int
	Retry    int
}

// NewJobQueue creates a new job queue with specified number of workers
func NewJobQueue(workers int, logger *slog.Logger) *JobQueue {
	ctx, cancel := context.WithCancel(context.Background())

	return &JobQueue{
		jobs:    make(chan Job, 100), // Buffered channel
		workers: workers,
		logger:  logger,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start starts the job queue workers
func (jq *JobQueue) Start() {
	for i := 0; i < jq.workers; i++ {
		jq.wg.Add(1)
		go jq.worker(i)
	}
}

// Stop stops the job queue and waits for all jobs to complete
func (jq *JobQueue) Stop() {
	close(jq.jobs)
	jq.cancel()
	jq.wg.Wait()
}

// worker runs in a goroutine and processes jobs
func (jq *JobQueue) worker(workerID int) {
	defer jq.wg.Done()

	for {
		select {
		case job, ok := <-jq.jobs:
			if !ok {
				return // Channel closed
			}

			jq.processJob(job)
		case <-jq.ctx.Done():
			return // Context cancelled
		}
	}
}

// processJob handles the execution of a single job with retry logic
func (jq *JobQueue) processJob(job Job) {
	var err error

	for job.Retry <= job.RetryMax {
		err = job.Task()
		if err == nil {
			// Success
			return
		}

		job.Retry++

		if job.Retry <= job.RetryMax {
			// Wait before retrying (exponential backoff)
			waitTime := time.Duration(job.Retry) * time.Second
			time.Sleep(waitTime)

			jq.logger.Info("Retrying job",
				slog.String("job_id", job.ID),
				slog.Int("retry_attempt", job.Retry),
				slog.String("error", err.Error()))
		}
	}

	// Log failure after all retries
	jq.logger.Error("Job failed after all retries",
		slog.String("job_id", job.ID),
		slog.String("error", err.Error()))
}

// Enqueue adds a job to the queue
func (jq *JobQueue) Enqueue(job Job) {
	select {
	case jq.jobs <- job:
		// Successfully added to queue
	default:
		// Queue is full, handle appropriately
		jq.logger.Error("Job queue is full, dropping job", slog.String("job_id", job.ID))
	}
}
