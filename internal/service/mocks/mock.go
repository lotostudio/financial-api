// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/lotostudio/financial-api/internal/domain"
)

// MockUsers is a mock of Users interface.
type MockUsers struct {
	ctrl     *gomock.Controller
	recorder *MockUsersMockRecorder
}

// MockUsersMockRecorder is the mock recorder for MockUsers.
type MockUsersMockRecorder struct {
	mock *MockUsers
}

// NewMockUsers creates a new mock instance.
func NewMockUsers(ctrl *gomock.Controller) *MockUsers {
	mock := &MockUsers{ctrl: ctrl}
	mock.recorder = &MockUsersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsers) EXPECT() *MockUsersMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockUsers) Get(ctx context.Context, userID int64) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, userID)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUsersMockRecorder) Get(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUsers)(nil).Get), ctx, userID)
}

// List mocks base method.
func (m *MockUsers) List(ctx context.Context) ([]domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockUsersMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockUsers)(nil).List), ctx)
}

// UpdatePassword mocks base method.
func (m *MockUsers) UpdatePassword(ctx context.Context, userID int64, toUpdate domain.UserToUpdate) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePassword", ctx, userID, toUpdate)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePassword indicates an expected call of UpdatePassword.
func (mr *MockUsersMockRecorder) UpdatePassword(ctx, userID, toUpdate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePassword", reflect.TypeOf((*MockUsers)(nil).UpdatePassword), ctx, userID, toUpdate)
}

// MockAuth is a mock of Auth interface.
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth.
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance.
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockAuth) Login(ctx context.Context, user domain.UserToLogin) (domain.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, user)
	ret0, _ := ret[0].(domain.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthMockRecorder) Login(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuth)(nil).Login), ctx, user)
}

// Refresh mocks base method.
func (m *MockAuth) Refresh(ctx context.Context, token string) (domain.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", ctx, token)
	ret0, _ := ret[0].(domain.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Refresh indicates an expected call of Refresh.
func (mr *MockAuthMockRecorder) Refresh(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockAuth)(nil).Refresh), ctx, token)
}

// Register mocks base method.
func (m *MockAuth) Register(ctx context.Context, user domain.UserToCreate) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, user)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockAuthMockRecorder) Register(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAuth)(nil).Register), ctx, user)
}

// MockCurrencies is a mock of Currencies interface.
type MockCurrencies struct {
	ctrl     *gomock.Controller
	recorder *MockCurrenciesMockRecorder
}

// MockCurrenciesMockRecorder is the mock recorder for MockCurrencies.
type MockCurrenciesMockRecorder struct {
	mock *MockCurrencies
}

// NewMockCurrencies creates a new mock instance.
func NewMockCurrencies(ctrl *gomock.Controller) *MockCurrencies {
	mock := &MockCurrencies{ctrl: ctrl}
	mock.recorder = &MockCurrenciesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCurrencies) EXPECT() *MockCurrenciesMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockCurrencies) List(ctx context.Context) ([]domain.Currency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]domain.Currency)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockCurrenciesMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockCurrencies)(nil).List), ctx)
}

// MockAccounts is a mock of Accounts interface.
type MockAccounts struct {
	ctrl     *gomock.Controller
	recorder *MockAccountsMockRecorder
}

// MockAccountsMockRecorder is the mock recorder for MockAccounts.
type MockAccountsMockRecorder struct {
	mock *MockAccounts
}

// NewMockAccounts creates a new mock instance.
func NewMockAccounts(ctrl *gomock.Controller) *MockAccounts {
	mock := &MockAccounts{ctrl: ctrl}
	mock.recorder = &MockAccountsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccounts) EXPECT() *MockAccountsMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAccounts) Create(ctx context.Context, toCreate domain.AccountToCreate, userID int64, currencyID int) (domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, toCreate, userID, currencyID)
	ret0, _ := ret[0].(domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockAccountsMockRecorder) Create(ctx, toCreate, userID, currencyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAccounts)(nil).Create), ctx, toCreate, userID, currencyID)
}

// Delete mocks base method.
func (m *MockAccounts) Delete(ctx context.Context, id, userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAccountsMockRecorder) Delete(ctx, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAccounts)(nil).Delete), ctx, id, userID)
}

// Get mocks base method.
func (m *MockAccounts) Get(ctx context.Context, id, userID int64) (domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id, userID)
	ret0, _ := ret[0].(domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockAccountsMockRecorder) Get(ctx, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAccounts)(nil).Get), ctx, id, userID)
}

// List mocks base method.
func (m *MockAccounts) List(ctx context.Context, userID int64) ([]domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, userID)
	ret0, _ := ret[0].([]domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockAccountsMockRecorder) List(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockAccounts)(nil).List), ctx, userID)
}

// ListGrouped mocks base method.
func (m *MockAccounts) ListGrouped(ctx context.Context, userID int64) (domain.GroupedAccounts, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListGrouped", ctx, userID)
	ret0, _ := ret[0].(domain.GroupedAccounts)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListGrouped indicates an expected call of ListGrouped.
func (mr *MockAccountsMockRecorder) ListGrouped(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListGrouped", reflect.TypeOf((*MockAccounts)(nil).ListGrouped), ctx, userID)
}

// Update mocks base method.
func (m *MockAccounts) Update(ctx context.Context, toUpdate domain.AccountToUpdate, id, userID int64) (domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, toUpdate, id, userID)
	ret0, _ := ret[0].(domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockAccountsMockRecorder) Update(ctx, toUpdate, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAccounts)(nil).Update), ctx, toUpdate, id, userID)
}

// MockAccountTypes is a mock of AccountTypes interface.
type MockAccountTypes struct {
	ctrl     *gomock.Controller
	recorder *MockAccountTypesMockRecorder
}

// MockAccountTypesMockRecorder is the mock recorder for MockAccountTypes.
type MockAccountTypesMockRecorder struct {
	mock *MockAccountTypes
}

// NewMockAccountTypes creates a new mock instance.
func NewMockAccountTypes(ctrl *gomock.Controller) *MockAccountTypes {
	mock := &MockAccountTypes{ctrl: ctrl}
	mock.recorder = &MockAccountTypesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountTypes) EXPECT() *MockAccountTypesMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockAccountTypes) List(ctx context.Context) ([]domain.AccountType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]domain.AccountType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockAccountTypesMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockAccountTypes)(nil).List), ctx)
}

// MockTransactions is a mock of Transactions interface.
type MockTransactions struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionsMockRecorder
}

// MockTransactionsMockRecorder is the mock recorder for MockTransactions.
type MockTransactionsMockRecorder struct {
	mock *MockTransactions
}

// NewMockTransactions creates a new mock instance.
func NewMockTransactions(ctrl *gomock.Controller) *MockTransactions {
	mock := &MockTransactions{ctrl: ctrl}
	mock.recorder = &MockTransactionsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactions) EXPECT() *MockTransactionsMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTransactions) Create(ctx context.Context, toCreate domain.TransactionToCreate, userID int64, categoryId, creditId, debitId *int64) (domain.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, toCreate, userID, categoryId, creditId, debitId)
	ret0, _ := ret[0].(domain.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTransactionsMockRecorder) Create(ctx, toCreate, userID, categoryId, creditId, debitId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTransactions)(nil).Create), ctx, toCreate, userID, categoryId, creditId, debitId)
}

// Delete mocks base method.
func (m *MockTransactions) Delete(ctx context.Context, id, userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTransactionsMockRecorder) Delete(ctx, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTransactions)(nil).Delete), ctx, id, userID)
}

// List mocks base method.
func (m *MockTransactions) List(ctx context.Context, filter domain.TransactionsFilter) ([]domain.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, filter)
	ret0, _ := ret[0].([]domain.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockTransactionsMockRecorder) List(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockTransactions)(nil).List), ctx, filter)
}

// MockTransactionCategories is a mock of TransactionCategories interface.
type MockTransactionCategories struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionCategoriesMockRecorder
}

// MockTransactionCategoriesMockRecorder is the mock recorder for MockTransactionCategories.
type MockTransactionCategoriesMockRecorder struct {
	mock *MockTransactionCategories
}

// NewMockTransactionCategories creates a new mock instance.
func NewMockTransactionCategories(ctrl *gomock.Controller) *MockTransactionCategories {
	mock := &MockTransactionCategories{ctrl: ctrl}
	mock.recorder = &MockTransactionCategoriesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionCategories) EXPECT() *MockTransactionCategoriesMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockTransactionCategories) List(ctx context.Context) ([]domain.TransactionCategory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]domain.TransactionCategory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockTransactionCategoriesMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockTransactionCategories)(nil).List), ctx)
}

// ListByType mocks base method.
func (m *MockTransactionCategories) ListByType(ctx context.Context, _type domain.TransactionType) ([]domain.TransactionCategory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByType", ctx, _type)
	ret0, _ := ret[0].([]domain.TransactionCategory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByType indicates an expected call of ListByType.
func (mr *MockTransactionCategoriesMockRecorder) ListByType(ctx, _type interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByType", reflect.TypeOf((*MockTransactionCategories)(nil).ListByType), ctx, _type)
}

// MockTransactionTypes is a mock of TransactionTypes interface.
type MockTransactionTypes struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionTypesMockRecorder
}

// MockTransactionTypesMockRecorder is the mock recorder for MockTransactionTypes.
type MockTransactionTypesMockRecorder struct {
	mock *MockTransactionTypes
}

// NewMockTransactionTypes creates a new mock instance.
func NewMockTransactionTypes(ctrl *gomock.Controller) *MockTransactionTypes {
	mock := &MockTransactionTypes{ctrl: ctrl}
	mock.recorder = &MockTransactionTypesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionTypes) EXPECT() *MockTransactionTypesMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockTransactionTypes) List(ctx context.Context) ([]domain.TransactionType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]domain.TransactionType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockTransactionTypesMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockTransactionTypes)(nil).List), ctx)
}
