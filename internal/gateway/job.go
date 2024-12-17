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

// CreateJobHandler creates a new job listing.
// @Summary Create a new job listing
// @Description Creates a new job listing with the provided job details. The job is associated with the authenticated user.
// @Tags Jobs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Job body service_models.Job true "Job Details"
// @Success 201 {object} service_models.Job "Job successfully created"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/jobs [post]
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
		internalServerError(w, r, err)
		return
	}

	if err = jsonResponse(w, http.StatusCreated, createdJob); err != nil {
		internalServerError(w, r, err)
	}

}

// GetAllJobsHandler retrieves all job listings.
// @Summary Retrieve all job listings
// @Description Fetches a list of all job listings available in the system.
// @Tags Jobs
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} service_models.Job "List of all jobs"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/jobs [get]
func (j *job) GetAllJobsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jobs, err := j.jobService.GetAllJobs(ctx)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	if err = jsonResponse(w, http.StatusOK, jobs); err != nil {
		internalServerError(w, r, err)
		return
	}
}

// GetAllJobsByUserHandler retrieves all job listings by a specific user.
// @Summary Retrieve all job listings by user ID
// @Description Fetches a list of all job listings associated with a specific user based on the user ID.
// @Tags Jobs
// @Produce json
// @Security ApiKeyAuth
// @Param userID path int64 true "User ID"
// @Success 200 {array} service_models.Job "List of jobs for the specified user"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/jobsByUser [get]
func (j *job) GetAllJobsByUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userID := r.Context().Value("userID").(int64)
	jobs, err := j.jobService.GetAllJobsByUserID(ctx, userID)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	if err = jsonResponse(w, http.StatusOK, jobs); err != nil {
		internalServerError(w, r, err)
		return
	}
}

// GetJobByIdHandler retrieves a specific job listing by its ID.
// @Summary Retrieve a job listing by ID
// @Description Fetches a job listing based on the provided job ID.
// @Tags Jobs
// @Produce json
// @Security ApiKeyAuth
// @Param id path int64 true "Job ID"
// @Success 200 {object} service_models.Job "Job listing details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/jobs/{id} [get]
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
		return
	}
}

// UpdateJobHandler updates an existing job listing.
// @Summary Update an existing job listing
// @Description Updates the details of a job listing, based on the provided job ID and job data.
// @Tags Jobs
// @Produce json
// @Security ApiKeyAuth
// @Param id path int64 true "Job ID"
// @Param job body service_models.Job true "Job data to update"
// @Success 200 {object} service_models.Job "Updated job details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/jobs/{id} [put]
func (j *job) UpdateJobHandler(w http.ResponseWriter, r *http.Request) {
	// TODO:Fix the issue of update in the update Job handler
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

// DeleteJobHandler deletes an existing job listing.
// @Summary Delete a job listing
// @Description Deletes a job listing by its ID. Only the user who created the job or an admin can delete it.
// @Tags Jobs
// @Produce json
// @Security ApiKeyAuth
// @Param id path int64 true "Job ID"
// @Success 200 {string} string "The job was successfully deleted"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/jobs/{id} [delete]
func (j *job) DeleteJobHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id, err := readIDParam(r)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	userID := r.Context().Value("userID").(int64)
	isAdmin := r.Context().Value("isAdmin").(bool)

	if err = j.jobService.DeleteJob(ctx, id, userID, isAdmin); err != nil {
		internalServerError(w, r, err)
		return
	}

	if err = jsonResponse(w, http.StatusOK, "the job was successfully deleted"); err != nil {
		internalServerError(w, r, err)
		return
	}
}

func NewJob(jobService service.Job) *job {
	return &job{
		jobService: jobService,
	}
}
