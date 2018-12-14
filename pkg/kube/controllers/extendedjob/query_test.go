package extendedjob_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/runtime"
	crc "sigs.k8s.io/controller-runtime/pkg/client"

	"code.cloudfoundry.org/cf-operator/pkg/kube/controllers/extendedjob"
	"code.cloudfoundry.org/cf-operator/pkg/kube/controllers/fakes"
	"code.cloudfoundry.org/cf-operator/testing"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// events := r.query.RecentPodEvents()
// pods := r.query.Match(events, extendedJob.Spec.Triggers.Selector)
// for _, pod := range pods {
//         if ok := r.matcher.Match(extendedJob, pod); ok {
//                 r.createJob()
//         }

var _ = Describe("Query", func() {
	var (
		client fakes.FakeClient
		env    testing.Catalog
	)

	Describe("RecentPodEvents", func() {

		act := func() ([]corev1.Event, error) {
			q := extendedjob.NewQuery(&client)
			return q.RecentPodEvents()
		}

		Context("when events exist", func() {
			BeforeEach(func() {
				now := time.Now()
				items := []corev1.Event{
					env.DefaultPodEvent(now.Add(extendedjob.Backlog - 1)),
					env.DefaultPodEvent(now.Add(-10 * time.Minute)),
					env.DefaultPodEvent(now.Add(-20 * time.Minute)),
				}
				listStub := func(ctx context.Context, ops *crc.ListOptions, obj runtime.Object) error {
					if list, ok := obj.(*corev1.EventList); ok {
						list.Items = items
					}
					return nil
				}
				client.ListCalls(listStub)
			})

			It("returns all pod related events with changes in specified timeframe", func() {
				events, err := act()
				Expect(err).ToNot(HaveOccurred())
				Expect(client.ListCallCount()).To(Equal(1))
				_, options, _ := client.ListArgsForCall(0)
				Expect(options.FieldSelector.String()).To(Equal("involvedObject.kind=Pod"))
				Expect(len(events)).To(Equal(2))
			})
		})

		It("returns error if client list fails", func() {
			client.ListReturns(fmt.Errorf("fake-error"))
			_, err := act()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("PodCandidates", func() {

		var events []corev1.Event

		act := func() ([]corev1.Pod, error) {
			q := extendedjob.NewQuery(&client)
			return q.PodCandidates(events)
		}

		Context("when events exist", func() {
			BeforeEach(func() {
				getStub := func(ctx context.Context, key types.NamespacedName, obj runtime.Object) error {
					if _, ok := obj.(*corev1.Pod); ok {
						obj = &corev1.Pod{}
					}
					return nil
				}
				client.GetCalls(getStub)
			})

			It("returns unique pods candidates for running jobs", func() {
				events = []corev1.Event{
					corev1.Event{InvolvedObject: corev1.ObjectReference{Name: "one"}},
					corev1.Event{InvolvedObject: corev1.ObjectReference{Name: "one"}},
					corev1.Event{InvolvedObject: corev1.ObjectReference{Name: "two"}},
				}

				pods, err := act()
				Expect(err).ToNot(HaveOccurred())
				Expect(len(pods)).To(Equal(2))
			})
		})

		// TODO
		Context("when events exist", func() {
			It("skips on if event references unknown pod", func() {})
		})

		It("returns error if get fails", func() {
			events = []corev1.Event{corev1.Event{}}
			client.GetReturns(fmt.Errorf("fake-error"))

			_, err := act()
			Expect(err).To(HaveOccurred())
		})
	})
})
