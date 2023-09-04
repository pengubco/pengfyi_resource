// Code generated by mockery v2.20.2. DO NOT EDIT.

package test_double_in_go

import mock "github.com/stretchr/testify/mock"

// MockGreeter is an autogenerated mock type for the Greeter type
type MockGreeter struct {
	mock.Mock
}

// Greeting provides a mock function with given fields: a
func (_m *MockGreeter) Greeting(a string) (string, error) {
	ret := _m.Called(a)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(a)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(a)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(a)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockGreeter interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockGreeter creates a new instance of MockGreeter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockGreeter(t mockConstructorTestingTNewMockGreeter) *MockGreeter {
	mock := &MockGreeter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}