package comprehend

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/pkg/errors"
)

type ComprehendJob struct {
	ctx    context.Context
	cancel context.CancelFunc
	svc    *comprehend.Comprehend
	jobID  string
}

func NewComprehendJob(
	ctx context.Context,
	svc *comprehend.Comprehend,
) *ComprehendJob {
	ctx, cancel := context.WithCancel(ctx)
	return &ComprehendJob{
		ctx:    ctx,
		cancel: cancel,
		svc:    svc,
	}
}

func (j *ComprehendJob) Cancel() error {
	j.cancel()
	if _, err := j.svc.StopEntitiesDetectionJob(
		&comprehend.StopEntitiesDetectionJobInput{
			JobId: aws.String(j.jobID),
		},
	); err != nil {
		return errors.Wrap(err, "StopEntitiesDetectionJob")
	}
	return nil
}

func (j *ComprehendJob) Context() context.Context {
	return j.ctx
}

func (j *ComprehendJob) Start() error {
	if j.jobID != "" {
		return errors.New("job already started")
	}
	resp, err := j.svc.StartEntitiesDetectionJobWithContext(
		j.ctx,
		&comprehend.StartEntitiesDetectionJobInput{},
	)
	if err != nil {
		return errors.Wrap(err, "StartEntitiesDetectionJobWithContext")
	}
	j.svc.DetectEntities(&comprehend.DetectEntitiesInput{})
	j.jobID = aws.StringValue(resp.JobId)
	switch aws.StringValue(resp.JobStatus) {
	case "SUBMITTED", "IN_PROGRESS":
		return nil
	default:
		return errors.Errorf("unexpected job status: %s", aws.StringValue(resp.JobStatus))
	}
}
