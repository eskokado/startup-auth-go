package mocks

// import (
// 	"github.com/eskokado/startup-auth-go/backend/pkg/domain/entity"
// 	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
// 	"github.com/stretchr/testify/mock"
// )

// type MockUserInterface struct {
// 	mock.Mock
// }

// func (m *MockUserInterface) WithName(newName vo.Name) (*entity.User, error) {
// 	args := m.Called(newName)
// 	return args.Get(0).(*entity.User), args.Error(1)
// }

// type MockUser struct {
// 	mock.Mock
// }

// func (m *MockUser) WithName(newName vo.Name) (*entity.User, error) {
// 	args := m.Called(newName)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*entity.User), args.Error(1)
// }

// func (m *MockUser) WithPasswordHash(passwordHash vo.PasswordHash) (*entity.User, error) {
// 	args := m.Called(passwordHash)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*entity.User), args.Error(1)
// }

// func (m *MockUser) ID() vo.ID {
// 	args := m.Called()
// 	return args.Get(0).(vo.ID)
// }

// func (m *MockUser) PasswordHash() vo.PasswordHash {
// 	args := m.Called()
// 	return args.Get(0).(vo.PasswordHash)
// }

// // Implemente outros métodos necessários da interface User
// func (m *MockUser) VerifyPassword(password string) bool {
// 	args := m.Called(password)
// 	return args.Bool(0)
// }
