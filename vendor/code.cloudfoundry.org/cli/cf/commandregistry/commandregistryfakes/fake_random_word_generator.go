// Code generated by counterfeiter. DO NOT EDIT.
package commandregistryfakes

import (
	"sync"

	"code.cloudfoundry.org/cli/cf/commandregistry"
)

type FakeRandomWordGenerator struct {
	BabbleStub        func() string
	babbleMutex       sync.RWMutex
	babbleArgsForCall []struct{}
	babbleReturns     struct {
		result1 string
	}
	babbleReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeRandomWordGenerator) Babble() string {
	fake.babbleMutex.Lock()
	ret, specificReturn := fake.babbleReturnsOnCall[len(fake.babbleArgsForCall)]
	fake.babbleArgsForCall = append(fake.babbleArgsForCall, struct{}{})
	fake.recordInvocation("Babble", []interface{}{})
	fake.babbleMutex.Unlock()
	if fake.BabbleStub != nil {
		return fake.BabbleStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.babbleReturns.result1
}

func (fake *FakeRandomWordGenerator) BabbleCallCount() int {
	fake.babbleMutex.RLock()
	defer fake.babbleMutex.RUnlock()
	return len(fake.babbleArgsForCall)
}

func (fake *FakeRandomWordGenerator) BabbleReturns(result1 string) {
	fake.BabbleStub = nil
	fake.babbleReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeRandomWordGenerator) BabbleReturnsOnCall(i int, result1 string) {
	fake.BabbleStub = nil
	if fake.babbleReturnsOnCall == nil {
		fake.babbleReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.babbleReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeRandomWordGenerator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.babbleMutex.RLock()
	defer fake.babbleMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeRandomWordGenerator) recordInvocation(key string, args []interface{}) {
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

var _ commandregistry.RandomWordGenerator = new(FakeRandomWordGenerator)
