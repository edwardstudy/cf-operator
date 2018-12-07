// +build !ignore_autogenerated

/*

Don't alter this file, it was generated.

*/
// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtendedJob) DeepCopyInto(out *ExtendedJob) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtendedJob.
func (in *ExtendedJob) DeepCopy() *ExtendedJob {
	if in == nil {
		return nil
	}
	out := new(ExtendedJob)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ExtendedJob) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtendedJobList) DeepCopyInto(out *ExtendedJobList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ExtendedJob, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtendedJobList.
func (in *ExtendedJobList) DeepCopy() *ExtendedJobList {
	if in == nil {
		return nil
	}
	out := new(ExtendedJobList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ExtendedJobList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtendedJobSpec) DeepCopyInto(out *ExtendedJobSpec) {
	*out = *in
	in.Triggers.DeepCopyInto(&out.Triggers)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtendedJobSpec.
func (in *ExtendedJobSpec) DeepCopy() *ExtendedJobSpec {
	if in == nil {
		return nil
	}
	out := new(ExtendedJobSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtendedJobStatus) DeepCopyInto(out *ExtendedJobStatus) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtendedJobStatus.
func (in *ExtendedJobStatus) DeepCopy() *ExtendedJobStatus {
	if in == nil {
		return nil
	}
	out := new(ExtendedJobStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in MatchLabels) DeepCopyInto(out *MatchLabels) {
	{
		in := &in
		*out = make(MatchLabels, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatchLabels.
func (in MatchLabels) DeepCopy() MatchLabels {
	if in == nil {
		return nil
	}
	out := new(MatchLabels)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Selector) DeepCopyInto(out *Selector) {
	*out = *in
	if in.MatchLabels != nil {
		in, out := &in.MatchLabels, &out.MatchLabels
		*out = make(MatchLabels, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Selector.
func (in *Selector) DeepCopy() *Selector {
	if in == nil {
		return nil
	}
	out := new(Selector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Triggers) DeepCopyInto(out *Triggers) {
	*out = *in
	in.Selector.DeepCopyInto(&out.Selector)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Triggers.
func (in *Triggers) DeepCopy() *Triggers {
	if in == nil {
		return nil
	}
	out := new(Triggers)
	in.DeepCopyInto(out)
	return out
}
