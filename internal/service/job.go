package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/saleh-ghazimoradi/GoJobs/internal/repository"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service/service_models"
)

type Job interface {
	CreateJob(ctx context.Context, job *service_models.Job) (*service_models.Job, error)
	GetAllJobs(ctx context.Context) ([]*service_models.Job, error)
	GetAllJobsByUserID(ctx context.Context, userID int64) ([]*service_models.Job, error)
	GetJobById(ctx context.Context, id int64) (*service_models.Job, error)
	UpdateJob(ctx context.Context, job *service_models.Job, userID int64, isAdmin bool) (*service_models.Job, error)
	GetWithTXT(tx *sql.Tx) Job
}

type jobService struct {
	jobRepo repository.Job
}

func (j *jobService) CreateJob(ctx context.Context, job *service_models.Job) (*service_models.Job, error) {
	return j.jobRepo.CreateJob(ctx, job)
}

func (j *jobService) GetAllJobs(ctx context.Context) ([]*service_models.Job, error) {
	return j.jobRepo.GetAllJobs(ctx)
}

func (j *jobService) GetAllJobsByUserID(ctx context.Context, userID int64) ([]*service_models.Job, error) {
	return j.jobRepo.GetAllJobsByUserID(ctx, userID)
}

func (j *jobService) GetJobById(ctx context.Context, id int64) (*service_models.Job, error) {
	return j.jobRepo.GetJobById(ctx, id)
}

func (j *jobService) UpdateJob(ctx context.Context, job *service_models.Job, userID int64, isAdmin bool) (*service_models.Job, error) {
	exisingJob, err := j.jobRepo.GetJobById(ctx, job.ID)
	if err != nil {
		return nil, err
	}

	if !isAdmin && exisingJob.UserID != userID {
		return nil, errors.New("not authorized")
	}
	return j.jobRepo.UpdateJob(ctx, job)
}

func (j *jobService) GetWithTXT(tx *sql.Tx) Job {
	return &jobService{
		jobRepo: j.jobRepo.GetWithTXT(tx),
	}
}

func NewJobService(jobRepo repository.Job) Job {
	return &jobService{
		jobRepo: jobRepo,
	}
}
