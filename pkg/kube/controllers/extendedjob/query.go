package extendedjob

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"code.cloudfoundry.org/cf-operator/pkg/kube/apis/extendedjob/v1alpha1"
)

// Backlog defines the maximal minutes passed for pod events we take into consideration
const Backlog = -30 * time.Minute

// Query for events involving pods and filter them
type Query interface {
	RecentPodEvents() ([]corev1.Event, error)
	PodCandidates([]corev1.Event) ([]corev1.Pod, error)
	// will need to check both, job timestamp and selector
	Match(v1alpha1.ExtendedJob, []corev1.Pod) []corev1.Pod
}

// NewQuery returns a new Query struct
func NewQuery(c client.Client) *QueryImpl {
	return &QueryImpl{client: c}
}

// QueryImpl implements the query interface
type QueryImpl struct {
	client client.Client
}

// RecentPodEvents returns all events involving pods from the past
func (q *QueryImpl) RecentPodEvents() ([]corev1.Event, error) {
	obj := &corev1.EventList{}
	sel := fields.Set{"involvedObject.kind": "Pod"}.AsSelector()
	err := q.client.List(context.TODO(), &client.ListOptions{FieldSelector: sel}, obj)
	if err != nil {
		return obj.Items, err
	}

	now := time.Now()
	items := []corev1.Event{}
	for _, ev := range obj.Items {
		if ev.LastTimestamp.Time.After(now.Add(Backlog)) {
			items = append(items, ev)
		}
	}

	return items, nil
}

func (q *QueryImpl) PodCandidates(events []corev1.Event) ([]corev1.Pod, error) {
	seen := map[types.NamespacedName]bool{}
	pods := []corev1.Pod{}
	for _, ev := range events {
		name := types.NamespacedName{Name: ev.InvolvedObject.Name, Namespace: ev.InvolvedObject.Namespace}
		if _, ok := seen[name]; !ok {
			seen[name] = true
			pod := &corev1.Pod{}
			err := q.client.Get(context.TODO(), name, pod)
			if err != nil {
				return pods, err
			}
			pods = append(pods, *pod)
		}
	}
	return pods, nil
}

func (q *QueryImpl) Match(job v1alpha1.ExtendedJob, events []corev1.Pod) []corev1.Pod {
	return []corev1.Pod{}
}
