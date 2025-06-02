package mocks

import "context"

type MockResetPasswordUseCase struct {
	Err error
}

func (m *MockResetPasswordUseCase) Execute(ctx context.Context, token string, password string) error {
	return m.Err
}
