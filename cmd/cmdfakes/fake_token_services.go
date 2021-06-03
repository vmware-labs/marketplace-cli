// Code generated by counterfeiter. DO NOT EDIT.
package cmdfakes

import (
	"sync"

	"github.com/vmware-labs/marketplace-cli/v2/cmd"
	"github.com/vmware-labs/marketplace-cli/v2/lib/csp"
)

type FakeTokenServices struct {
	RedeemStub        func(string) (*csp.Claims, error)
	redeemMutex       sync.RWMutex
	redeemArgsForCall []struct {
		arg1 string
	}
	redeemReturns struct {
		result1 *csp.Claims
		result2 error
	}
	redeemReturnsOnCall map[int]struct {
		result1 *csp.Claims
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeTokenServices) Redeem(arg1 string) (*csp.Claims, error) {
	fake.redeemMutex.Lock()
	ret, specificReturn := fake.redeemReturnsOnCall[len(fake.redeemArgsForCall)]
	fake.redeemArgsForCall = append(fake.redeemArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Redeem", []interface{}{arg1})
	fake.redeemMutex.Unlock()
	if fake.RedeemStub != nil {
		return fake.RedeemStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.redeemReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeTokenServices) RedeemCallCount() int {
	fake.redeemMutex.RLock()
	defer fake.redeemMutex.RUnlock()
	return len(fake.redeemArgsForCall)
}

func (fake *FakeTokenServices) RedeemCalls(stub func(string) (*csp.Claims, error)) {
	fake.redeemMutex.Lock()
	defer fake.redeemMutex.Unlock()
	fake.RedeemStub = stub
}

func (fake *FakeTokenServices) RedeemArgsForCall(i int) string {
	fake.redeemMutex.RLock()
	defer fake.redeemMutex.RUnlock()
	argsForCall := fake.redeemArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTokenServices) RedeemReturns(result1 *csp.Claims, result2 error) {
	fake.redeemMutex.Lock()
	defer fake.redeemMutex.Unlock()
	fake.RedeemStub = nil
	fake.redeemReturns = struct {
		result1 *csp.Claims
		result2 error
	}{result1, result2}
}

func (fake *FakeTokenServices) RedeemReturnsOnCall(i int, result1 *csp.Claims, result2 error) {
	fake.redeemMutex.Lock()
	defer fake.redeemMutex.Unlock()
	fake.RedeemStub = nil
	if fake.redeemReturnsOnCall == nil {
		fake.redeemReturnsOnCall = make(map[int]struct {
			result1 *csp.Claims
			result2 error
		})
	}
	fake.redeemReturnsOnCall[i] = struct {
		result1 *csp.Claims
		result2 error
	}{result1, result2}
}

func (fake *FakeTokenServices) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.redeemMutex.RLock()
	defer fake.redeemMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeTokenServices) recordInvocation(key string, args []interface{}) {
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

var _ cmd.TokenServices = new(FakeTokenServices)
