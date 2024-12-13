package gateway

import (
	"context"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service/service_models"
	"net/http"
	"time"
)

type job struct {
	jobService service.Job
}

func (j *job) CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var jobs service_models.Job
	if err := readJSON(w, r, &jobs); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(jobs); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	userID := r.Context().Value("userID").(int64)
	jobs.UserID = userID

	createdJob, err := j.jobService.CreateJob(ctx, &jobs)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, createdJob); err != nil {
		internalServerError(w, r, err)
	}

}

func NewJob(jobService service.Job) *job {
	return &job{
		jobService: jobService,
	}
}
