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

	if err = jsonResponse(w, http.StatusCreated, createdJob); err != nil {
		internalServerError(w, r, err)
	}

}

func (j *job) GetAllJobsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jobs, err := j.jobService.GetAllJobs(ctx)
	if err != nil {
		internalServerError(w, r, err)
	}

	if err = jsonResponse(w, http.StatusOK, jobs); err != nil {
		internalServerError(w, r, err)
	}
}

func (j *job) GetAllJobsByUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userID := r.Context().Value("userID").(int64)
	jobs, err := j.jobService.GetAllJobsByUserID(ctx, userID)
	if err != nil {
		internalServerError(w, r, err)
	}

	if err = jsonResponse(w, http.StatusOK, jobs); err != nil {
		internalServerError(w, r, err)
	}
}

func (j *job) GetJobByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id, err := readIDParam(r)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	jobs, err := j.jobService.GetJobById(ctx, id)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	if err = jsonResponse(w, http.StatusOK, jobs); err != nil {
		internalServerError(w, r, err)
	}
}

func (j *job) UpdateJobHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id, err := readIDParam(r)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	var jobs service_models.Job
	jobs.ID = id

	if err = readJSON(w, r, &jobs); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err = Validate.Struct(jobs); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	userID := r.Context().Value("userID").(int64)
	isAdmin := r.Context().Value("isAdmin").(bool)

	updateJob, err := j.jobService.UpdateJob(ctx, &jobs, userID, isAdmin)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	if err = jsonResponse(w, http.StatusOK, updateJob); err != nil {
		internalServerError(w, r, err)
		return
	}

}

func NewJob(jobService service.Job) *job {
	return &job{
		jobService: jobService,
	}
}
