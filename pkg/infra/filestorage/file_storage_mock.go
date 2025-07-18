// Code generated by mockery v2.53.4. DO NOT EDIT.

package filestorage

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockFileStorage is an autogenerated mock type for the FileStorage type
type MockFileStorage struct {
	mock.Mock
}

// CreateFolder provides a mock function with given fields: ctx, path
func (_m *MockFileStorage) CreateFolder(ctx context.Context, path string) error {
	ret := _m.Called(ctx, path)

	if len(ret) == 0 {
		panic("no return value specified for CreateFolder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, path
func (_m *MockFileStorage) Delete(ctx context.Context, path string) error {
	ret := _m.Called(ctx, path)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteFolder provides a mock function with given fields: ctx, path, options
func (_m *MockFileStorage) DeleteFolder(ctx context.Context, path string, options *DeleteFolderOptions) error {
	ret := _m.Called(ctx, path, options)

	if len(ret) == 0 {
		panic("no return value specified for DeleteFolder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *DeleteFolderOptions) error); ok {
		r0 = rf(ctx, path, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, path, options
func (_m *MockFileStorage) Get(ctx context.Context, path string, options *GetFileOptions) (*File, bool, error) {
	ret := _m.Called(ctx, path, options)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *File
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *GetFileOptions) (*File, bool, error)); ok {
		return rf(ctx, path, options)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *GetFileOptions) *File); ok {
		r0 = rf(ctx, path, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*File)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *GetFileOptions) bool); ok {
		r1 = rf(ctx, path, options)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, *GetFileOptions) error); ok {
		r2 = rf(ctx, path, options)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// List provides a mock function with given fields: ctx, folderPath, paging, options
func (_m *MockFileStorage) List(ctx context.Context, folderPath string, paging *Paging, options *ListOptions) (*ListResponse, error) {
	ret := _m.Called(ctx, folderPath, paging, options)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 *ListResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *Paging, *ListOptions) (*ListResponse, error)); ok {
		return rf(ctx, folderPath, paging, options)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *Paging, *ListOptions) *ListResponse); ok {
		r0 = rf(ctx, folderPath, paging, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ListResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *Paging, *ListOptions) error); ok {
		r1 = rf(ctx, folderPath, paging, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Upsert provides a mock function with given fields: ctx, command
func (_m *MockFileStorage) Upsert(ctx context.Context, command *UpsertFileCommand) error {
	ret := _m.Called(ctx, command)

	if len(ret) == 0 {
		panic("no return value specified for Upsert")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *UpsertFileCommand) error); ok {
		r0 = rf(ctx, command)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// close provides a mock function with no fields
func (_m *MockFileStorage) close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockFileStorage creates a new instance of MockFileStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockFileStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockFileStorage {
	mock := &MockFileStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
