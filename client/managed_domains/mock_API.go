// Code generated by mockery v1.0.0. DO NOT EDIT.

package managed_domains

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockAPI is an autogenerated mock type for the API type
type MockAPI struct {
	mock.Mock
}

// ListManagedDomains provides a mock function with given fields: ctx, params
func (_m *MockAPI) ListManagedDomains(ctx context.Context, params *ListManagedDomainsParams) (*ListManagedDomainsOK, error) {
	ret := _m.Called(ctx, params)

	var r0 *ListManagedDomainsOK
	if rf, ok := ret.Get(0).(func(context.Context, *ListManagedDomainsParams) *ListManagedDomainsOK); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ListManagedDomainsOK)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ListManagedDomainsParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
