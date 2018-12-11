package extendedjob

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"code.cloudfoundry.org/cf-operator/pkg/kube/apis/extendedjob/v1alpha1"
)

type Runner interface {
	Run()
}

func NewRunner(
	log *zap.SugaredLogger,
	mgr manager.Manager,
) *RunnerImpl {
	return &RunnerImpl{
		log:      log,
		client:   mgr.GetClient(),
		recorder: mgr.GetRecorder("extendedjob runner"),
	}
}

type RunnerImpl struct {
	log      *zap.SugaredLogger
	client   client.Client
	recorder record.EventRecorder
}

// Run checks all existing extendedJobs
//&JobRunner{}
// get all exJobs
// filter pod events
// run jobs
//   only once per event
//   only for latest
//   exJob is owner of Job
func (r *RunnerImpl) Run() {
	obj := &v1alpha1.ExtendedJobList{}
	err := r.client.List(context.TODO(), &client.ListOptions{}, obj)
	if err != nil {
		r.log.Infof("failed to query extended jobs: %s", err)
		return
	}

	for i := 0; i < len(obj.Items); i++ {
		extendedJob := obj.Items[i]

		job := &batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("job-%s-%d", extendedJob.Name, i),
				Namespace: extendedJob.Namespace,
			},
			Spec: jobSpec(),
		}

		err = r.client.Create(context.TODO(), job)
		if err != nil {
			r.log.Infof("failed to create job for %s: %s", extendedJob.Name, err)
		}

	}
}

func jobSpec() batchv1.JobSpec {
	one := int64(1)
	return batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				TerminationGracePeriodSeconds: &one,
				Containers: []corev1.Container{
					{
						Name:    "busybox",
						Image:   "busybox",
						Command: []string{"sleep", "6"},
					},
				},
			},
		},
	}

}
