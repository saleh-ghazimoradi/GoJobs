package service

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/GoJobs/internal/repository"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service/service_models"
)

type Job interface {
	CreateJob(ctx context.Context, job *service_models.Job) (*service_models.Job, error)
	GetAllJobs(ctx context.Context) ([]*service_models.Job, error)
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
