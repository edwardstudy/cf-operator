// Code generated by counterfeiter. DO NOT EDIT.
package manifestfakes

import (
	sync "sync"

	v1alpha1 "code.cloudfoundry.org/cf-operator/pkg/apis/fissile/v1alpha1"
	manifest "code.cloudfoundry.org/cf-operator/pkg/bosh/manifest"
)

type FakeResolver struct {
	ResolveCRDStub        func(v1alpha1.BOSHDeploymentSpec, string) (*manifest.Manifest, error)
	resolveCRDMutex       sync.RWMutex
	resolveCRDArgsForCall []struct {
		arg1 v1alpha1.BOSHDeploymentSpec
		arg2 string
	}
	resolveCRDReturns struct {
		result1 *manifest.Manifest
		result2 error
	}
	resolveCRDReturnsOnCall map[int]struct {
		result1 *manifest.Manifest
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeResolver) ResolveCRD(arg1 v1alpha1.BOSHDeploymentSpec, arg2 string) (*manifest.Manifest, error) {
	fake.resolveCRDMutex.Lock()
	ret, specificReturn := fake.resolveCRDReturnsOnCall[len(fake.resolveCRDArgsForCall)]
	fake.resolveCRDArgsForCall = append(fake.resolveCRDArgsForCall, struct {
		arg1 v1alpha1.BOSHDeploymentSpec
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("ResolveCRD", []interface{}{arg1, arg2})
	fake.resolveCRDMutex.Unlock()
	if fake.ResolveCRDStub != nil {
		return fake.ResolveCRDStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.resolveCRDReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeResolver) ResolveCRDCallCount() int {
	fake.resolveCRDMutex.RLock()
	defer fake.resolveCRDMutex.RUnlock()
	return len(fake.resolveCRDArgsForCall)
}

func (fake *FakeResolver) ResolveCRDCalls(stub func(v1alpha1.BOSHDeploymentSpec, string) (*manifest.Manifest, error)) {
	fake.resolveCRDMutex.Lock()
	defer fake.resolveCRDMutex.Unlock()
	fake.ResolveCRDStub = stub
}

func (fake *FakeResolver) ResolveCRDArgsForCall(i int) (v1alpha1.BOSHDeploymentSpec, string) {
	fake.resolveCRDMutex.RLock()
	defer fake.resolveCRDMutex.RUnlock()
	argsForCall := fake.resolveCRDArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeResolver) ResolveCRDReturns(result1 *manifest.Manifest, result2 error) {
	fake.resolveCRDMutex.Lock()
	defer fake.resolveCRDMutex.Unlock()
	fake.ResolveCRDStub = nil
	fake.resolveCRDReturns = struct {
		result1 *manifest.Manifest
		result2 error
	}{result1, result2}
}

func (fake *FakeResolver) ResolveCRDReturnsOnCall(i int, result1 *manifest.Manifest, result2 error) {
	fake.resolveCRDMutex.Lock()
	defer fake.resolveCRDMutex.Unlock()
	fake.ResolveCRDStub = nil
	if fake.resolveCRDReturnsOnCall == nil {
		fake.resolveCRDReturnsOnCall = make(map[int]struct {
			result1 *manifest.Manifest
			result2 error
		})
	}
	fake.resolveCRDReturnsOnCall[i] = struct {
		result1 *manifest.Manifest
		result2 error
	}{result1, result2}
}

func (fake *FakeResolver) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.resolveCRDMutex.RLock()
	defer fake.resolveCRDMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeResolver) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ manifest.Resolver = new(FakeResolver)
