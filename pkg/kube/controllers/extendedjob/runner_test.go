package extendedjob_test

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/cf-operator/pkg/kube/apis/extendedjob/v1alpha1"
	"code.cloudfoundry.org/cf-operator/pkg/kube/controllers"
	"code.cloudfoundry.org/cf-operator/pkg/kube/controllers/extendedjob"
	"code.cloudfoundry.org/cf-operator/pkg/kube/controllers/fakes"
	"code.cloudfoundry.org/cf-operator/testing"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"k8s.io/client-go/kubernetes/scheme"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	crc "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Runner", func() {
	Describe("Run", func() {
		var (
			env    testing.Catalog
			mgr    *fakes.FakeManager
			logs   *observer.ObservedLogs
			log    *zap.SugaredLogger
			runner *extendedjob.RunnerImpl
		)

		BeforeEach(func() {
			controllers.AddToScheme(scheme.Scheme)
			logs, log = testing.NewTestLogger()
			mgr = &fakes.FakeManager{}
		})

		Context("when client fails", func() {
			var client fakes.FakeClient

			BeforeEach(func() {
				client = fakes.FakeClient{}
				mgr.GetClientReturns(&client)
			})

			It("should log list failure", func() {
				client.ListReturns(fmt.Errorf("fake-error"))

				extendedjob.NewRunner(log, mgr).Run()
				Expect(logs.FilterMessage("failed to query extended jobs: fake-error").Len()).To(Equal(1))
			})

			It("should log create error and continue", func() {
				jobList := []v1alpha1.ExtendedJob{
					*env.DefaultExtendedJob("foo"),
					*env.DefaultExtendedJob("bar"),
				}
				listStub := func(ctx context.Context, ops *crc.ListOptions, obj runtime.Object) error {
					if list, ok := obj.(*v1alpha1.ExtendedJobList); ok {
						list.Items = jobList
					}
					return nil
				}
				client.ListCalls(listStub)
				client.CreateReturns(fmt.Errorf("fake-error"))

				extendedjob.NewRunner(log, mgr).Run()
				Expect(client.CreateCallCount()).To(Equal(2))
				Expect(logs.FilterMessageSnippet("failed to create job for foo: fake-error").Len()).To(Equal(1))
				Expect(logs.FilterMessageSnippet("failed to create job for bar: fake-error").Len()).To(Equal(1))
			})
		})

		Context("when extended jobs are present", func() {
			var (
				exJobs []runtime.Object
				client crc.Client
			)

			act := func() {
				client = fake.NewFakeClient(exJobs...)
				mgr.GetClientReturns(client)

				runner = extendedjob.NewRunner(log, mgr)
				runner.Run()
			}

			Context("when no extended job is present", func() {
				BeforeEach(func() {
					exJobs = []runtime.Object{}
				})

				It("should not create jobs", func() {
					act()

					obj := &batchv1.JobList{}
					err := client.List(context.TODO(), &crc.ListOptions{}, obj)
					Expect(err).ToNot(HaveOccurred())
					Expect(len(obj.Items)).To(Equal(0))
				})
			})

			Context("when extended jobs are present", func() {
				BeforeEach(func() {
					exJobs = []runtime.Object{env.DefaultExtendedJob("foo"), env.DefaultExtendedJob("bar")}
				})

				It("should create a job for each pod matched by a extendedjob", func() {
					act()

					obj := &batchv1.JobList{}
					err := client.List(context.TODO(), &crc.ListOptions{}, obj)
					Expect(err).ToNot(HaveOccurred())
					Expect(len(obj.Items)).To(Equal(2))
					Expect(obj.Items[0].Name).To(ContainSubstring("job-foo-"))
				})
			})

		})
	})
})

func listExtendedJobs(client client.Client) (*v1alpha1.ExtendedJobList, error) {
	obj := &v1alpha1.ExtendedJobList{}
	err := client.List(
		context.TODO(),
		&crc.ListOptions{
			Raw: &metav1.ListOptions{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ExtendedJob",
					APIVersion: "fissile.cloudfoundry.org/v1alpha1",
				},
			},
		},
		obj,
	)
	return obj, err
}
