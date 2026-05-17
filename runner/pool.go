package runner

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

type JobStatus string

const (
	JobQueued    JobStatus = "queued"
	JobRunning   JobStatus = "running"
	JobCompleted JobStatus = "completed"
	JobFailed    JobStatus = "failed"
)

type JobRecord struct {
	Status JobStatus
	Output string
	Error  string
}

type CompileJob struct {
	ID         string
	SourceCode string
}

var (
	JobQueue chan CompileJob
	jobsMu   sync.RWMutex
	jobs     = make(map[string]*JobRecord)
)

func InitWorkerPool(maxWorkers, maxQueueSize int) {
	JobQueue = make(chan CompileJob, maxQueueSize) // buffered chan

	for i := 0; i < maxWorkers; i++ {
		go worker(i, JobQueue) // go-routine
	}

	log.Printf("wrker pool started with %d workers", maxWorkers)
}

func EnqueueJob(sourceCode string) string {
	jobID := uuid.New().String()

	jobsMu.Lock()
	jobs[jobID] = &JobRecord{Status: JobQueued}
	jobsMu.Unlock()

	JobQueue <- CompileJob{ID: jobID, SourceCode: sourceCode}

	return jobID
}

func GetJob(jobID string) (JobRecord, bool) {
	jobsMu.RLock()
	record, ok := jobs[jobID]
	if !ok {
		jobsMu.RUnlock()
		return JobRecord{}, false
	}
	snapshot := *record
	jobsMu.RUnlock()

	return snapshot, true
}

func updateJob(jobID string, status JobStatus, output string, errMsg string) {
	jobsMu.Lock()
	record, ok := jobs[jobID]
	if !ok {
		record = &JobRecord{}
		jobs[jobID] = record
	}
	record.Status = status
	record.Output = output
	record.Error = errMsg
	jobsMu.Unlock()
}

// job
func worker(id int, jobs <-chan CompileJob) {
	for job := range jobs {
		log.Printf("worker %d started working...", id)

		updateJob(job.ID, JobRunning, "", "")
		output, err := RunSource(job.SourceCode)
		if err != nil {
			updateJob(job.ID, JobFailed, "", err.Error())
		} else {
			updateJob(job.ID, JobCompleted, output, "")
		}

		log.Printf("worker %d completed the task...", id)
	}
}