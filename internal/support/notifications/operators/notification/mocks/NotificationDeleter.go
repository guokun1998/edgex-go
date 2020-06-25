// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// NotificationDeleter is an autogenerated mock type for the NotificationDeleter type
type NotificationDeleter struct {
	mock.Mock
}

// DeleteNotificationById provides a mock function with given fields: id
func (_m *NotificationDeleter) DeleteNotificationById(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteNotificationBySlug provides a mock function with given fields: slug
func (_m *NotificationDeleter) DeleteNotificationBySlug(slug string) error {
	ret := _m.Called(slug)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(slug)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteNotificationsOld provides a mock function with given fields: age
func (_m *NotificationDeleter) DeleteNotificationsOld(age int) error {
	ret := _m.Called(age)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(age)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}