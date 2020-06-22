package job

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"strings"

	logutil "github.com/filanov/bm-inventory/pkg/log"
	"github.com/sirupsen/logrus"
	batch "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type localJob struct {
	Config
	log logrus.FieldLogger
}

func NewLocalJob(log logrus.FieldLogger, cfg Config) *localJob {
	return &localJob{
		Config: cfg,
		log:    log,
	}
}

func (j *localJob) Create(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error {
	log := logutil.FromContext(ctx, j.log)
	job := obj.(*batch.Job)
	if job.TypeMeta.Kind != "Job" {
		// Skip if not a Job
		return nil
	}

	jobName := job.ObjectMeta.Name

	if strings.HasPrefix(jobName, "generate-kubeconfig") {
		log.Infof("Executing generate-kubeconfig locally for job %s", jobName)
		return j.executeKubeconfigJob(job)
	}

	if strings.HasPrefix(jobName, "createimage") || strings.HasPrefix(jobName, "dummyimage") {
		log.Infof("Creating ISO locally for job %s", jobName)
		return j.executeImageJob(job)
	}

	return nil
}

func (j *localJob) executeKubeconfigJob(job *batch.Job) error {
	cmd := exec.Command("python", "./data/process-ignition-manifests-and-kubeconfig.py")
	cmd.Env = append(os.Environ(),
		"S3_ENDPOINT_URL="+job.Spec.Template.Spec.Containers[0].Env[0].Value,
		"INSTALLER_CONFIG="+job.Spec.Template.Spec.Containers[0].Env[1].Value,
		"IMAGE_NAME="+job.Spec.Template.Spec.Containers[0].Env[2].Value,
		"S3_BUCKET="+job.Spec.Template.Spec.Containers[0].Env[3].Value,
		"CLUSTER_ID="+job.Spec.Template.Spec.Containers[0].Env[4].Value,
		"OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE="+job.Spec.Template.Spec.Containers[0].Env[5].Value,
		"aws_access_key_id="+job.Spec.Template.Spec.Containers[0].Env[6].Value,
		"aws_secret_access_key="+job.Spec.Template.Spec.Containers[0].Env[7].Value,
		"WORK_DIR=/data",
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (j *localJob) executeImageJob(job *batch.Job) error {
	// TBD with MGMT-1175
	return nil
}

func (j *localJob) Monitor(ctx context.Context, name, namespace string) error {
	log := logutil.FromContext(ctx, j.log)
	log.Info("localJob.Monitor is NOOP")
	return nil
}

func (j *localJob) Delete(ctx context.Context, name, namespace string) error {
	log := logutil.FromContext(ctx, j.log)
	log.Info("localJob.Delete is NOOP")
	return nil
}
