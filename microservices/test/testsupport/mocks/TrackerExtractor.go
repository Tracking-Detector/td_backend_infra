// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// TrackerExtractor is an autogenerated mock type for the TrackerExtractor type
type TrackerExtractor struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *TrackerExtractor) Execute(_a0 bool) ([]int, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 []int
	var r1 error
	if rf, ok := ret.Get(0).(func(bool) ([]int, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(bool) []int); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	if rf, ok := ret.Get(1).(func(bool) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTrackerExtractor creates a new instance of TrackerExtractor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTrackerExtractor(t interface {
	mock.TestingT
	Cleanup(func())
}) *TrackerExtractor {
	mock := &TrackerExtractor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
