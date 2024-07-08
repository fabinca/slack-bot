// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	gojenkins "github.com/bndr/gojenkins"
	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// BuildJob provides a mock function with given fields: ctx, name, params
func (_m *Client) BuildJob(ctx context.Context, name string, params map[string]string) (int64, error) {
	ret := _m.Called(ctx, name, params)

	if len(ret) == 0 {
		panic("no return value specified for BuildJob")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]string) (int64, error)); ok {
		return rf(ctx, name, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]string) int64); ok {
		r0 = rf(ctx, name, params)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, map[string]string) error); ok {
		r1 = rf(ctx, name, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllNodes provides a mock function with given fields: ctx
func (_m *Client) GetAllNodes(ctx context.Context) ([]*gojenkins.Node, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllNodes")
	}

	var r0 []*gojenkins.Node
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*gojenkins.Node, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*gojenkins.Node); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*gojenkins.Node)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetJob provides a mock function with given fields: ctx, id
func (_m *Client) GetJob(ctx context.Context, id string) (*gojenkins.Job, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetJob")
	}

	var r0 *gojenkins.Job
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*gojenkins.Job, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *gojenkins.Job); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gojenkins.Job)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
