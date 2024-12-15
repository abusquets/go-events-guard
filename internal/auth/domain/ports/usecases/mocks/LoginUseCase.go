// Code generated by mockery v2.50.0. DO NOT EDIT.

package mocks

import (
	context "context"
	constants "eventsguard/internal/auth/constants"

	dtos "eventsguard/internal/auth/dtos"

	entities "eventsguard/internal/auth/domain/entities"

	errors "eventsguard/internal/app/errors"

	mock "github.com/stretchr/testify/mock"
)

// MockLoginUseCase is an autogenerated mock type for the LoginUseCase type
type MockLoginUseCase struct {
	mock.Mock
}

type MockLoginUseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLoginUseCase) EXPECT() *MockLoginUseCase_Expecter {
	return &MockLoginUseCase_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx, data, device
func (_m *MockLoginUseCase) Execute(ctx context.Context, data dtos.LoginInputDTO, device constants.TokenDevice) (*entities.Token, *errors.AppError) {
	ret := _m.Called(ctx, data, device)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 *entities.Token
	var r1 *errors.AppError
	if rf, ok := ret.Get(0).(func(context.Context, dtos.LoginInputDTO, constants.TokenDevice) (*entities.Token, *errors.AppError)); ok {
		return rf(ctx, data, device)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dtos.LoginInputDTO, constants.TokenDevice) *entities.Token); ok {
		r0 = rf(ctx, data, device)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Token)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, dtos.LoginInputDTO, constants.TokenDevice) *errors.AppError); ok {
		r1 = rf(ctx, data, device)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*errors.AppError)
		}
	}

	return r0, r1
}

// MockLoginUseCase_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockLoginUseCase_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - data dtos.LoginInputDTO
//   - device constants.TokenDevice
func (_e *MockLoginUseCase_Expecter) Execute(ctx interface{}, data interface{}, device interface{}) *MockLoginUseCase_Execute_Call {
	return &MockLoginUseCase_Execute_Call{Call: _e.mock.On("Execute", ctx, data, device)}
}

func (_c *MockLoginUseCase_Execute_Call) Run(run func(ctx context.Context, data dtos.LoginInputDTO, device constants.TokenDevice)) *MockLoginUseCase_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dtos.LoginInputDTO), args[2].(constants.TokenDevice))
	})
	return _c
}

func (_c *MockLoginUseCase_Execute_Call) Return(_a0 *entities.Token, _a1 *errors.AppError) *MockLoginUseCase_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockLoginUseCase_Execute_Call) RunAndReturn(run func(context.Context, dtos.LoginInputDTO, constants.TokenDevice) (*entities.Token, *errors.AppError)) *MockLoginUseCase_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockLoginUseCase creates a new instance of MockLoginUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockLoginUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockLoginUseCase {
	mock := &MockLoginUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}