package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service/service_models"
)

type Job interface {
	CreateJob(ctx context.Context, job *service_models.Job) (*service_models.Job, error)
	GetAllJobs(ctx context.Context) ([]*service_models.Job, error)
	GetWithTXT(tx *sql.Tx) Job
}

type jobRepository struct {
	dbWrite *sql.DB
	dbRead  *sql.DB
	tx      *sql.Tx
}

func (j *jobRepository) CreateJob(ctx context.Context, job *service_models.Job) (*service_models.Job, error) {
	query := `INSERT INTO jobs (title, description, company, location, salary,user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at;`
	var id int64
	err := j.dbWrite.QueryRowContext(ctx, query, job.Title, job.Description, job.Company, job.Location, job.Salary, job.UserID).Scan(&id, &job.CreatedAt)
	if err != nil {
		return nil, err
	}
	job.ID = id
	return job, nil
}

func (j *jobRepository) GetAllJobs(ctx context.Context) ([]*service_models.Job, error) {
	query := `SELECT id, title, description, location, company, salary, created_at, user_id FROM jobs`
	rows, err := j.dbRead.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []*service_models.Job
	for rows.Next() {
		var job service_models.Job
		if err = rows.Scan(&job.ID, &job.Title, &job.Description, &job.Location, &job.Company, &job.Salary, &job.CreatedAt, &job.UserID); err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (j *jobRepository) GetWithTXT(tx *sql.Tx) Job {
	return &jobRepository{
		dbWrite: j.dbWrite,
		dbRead:  j.dbRead,
		tx:      tx,
	}
}

func NewJobRepository(dbWrite *sql.DB, dbRead *sql.DB) Job {
	return &jobRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
