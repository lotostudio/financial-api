// Code generated by MockGen. DO NOT EDIT.
// Source: repo.go

// Package mock_repo is a generated GoMock package.
package mock_repo

import (
	context "context"
	reflect "reflect"
	time "time"

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

// Create mocks base method.
func (m *MockUsers) Create(ctx context.Context, user domain.User) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUsersMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUsers)(nil).Create), ctx, user)
}

// Get mocks base method.
func (m *MockUsers) Get(ctx context.Context, id int64) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUsersMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUsers)(nil).Get), ctx, id)
}

// GetByCredentials mocks base method.
func (m *MockUsers) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByCredentials", ctx, email, password)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByCredentials indicates an expected call of GetByCredentials.
func (mr *MockUsersMockRecorder) GetByCredentials(ctx, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByCredentials", reflect.TypeOf((*MockUsers)(nil).GetByCredentials), ctx, email, password)
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

// MockSessions is a mock of Sessions interface.
type MockSessions struct {
	ctrl     *gomock.Controller
	recorder *MockSessionsMockRecorder
}

// MockSessionsMockRecorder is the mock recorder for MockSessions.
type MockSessionsMockRecorder struct {
	mock *MockSessions
}

// NewMockSessions creates a new mock instance.
func NewMockSessions(ctrl *gomock.Controller) *MockSessions {
	mock := &MockSessions{ctrl: ctrl}
	mock.recorder = &MockSessionsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessions) EXPECT() *MockSessionsMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSessions) Create(ctx context.Context, userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockSessionsMockRecorder) Create(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSessions)(nil).Create), ctx, userID)
}

// GetByToken mocks base method.
func (m *MockSessions) GetByToken(ctx context.Context, token string) (domain.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByToken", ctx, token)
	ret0, _ := ret[0].(domain.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByToken indicates an expected call of GetByToken.
func (mr *MockSessionsMockRecorder) GetByToken(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByToken", reflect.TypeOf((*MockSessions)(nil).GetByToken), ctx, token)
}

// Update mocks base method.
func (m *MockSessions) Update(ctx context.Context, toUpdate domain.SessionToUpdate, userID int64) (domain.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, toUpdate, userID)
	ret0, _ := ret[0].(domain.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockSessionsMockRecorder) Update(ctx, toUpdate, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSessions)(nil).Update), ctx, toUpdate, userID)
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

// Get mocks base method.
func (m *MockCurrencies) Get(ctx context.Context, id int) (domain.Currency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(domain.Currency)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCurrenciesMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCurrencies)(nil).Get), ctx, id)
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

// CountByTypes mocks base method.
func (m *MockAccounts) CountByTypes(ctx context.Context, userID int64, _type domain.AccountType, types ...domain.AccountType) (int64, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, userID, _type}
	for _, a := range types {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CountByTypes", varargs...)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountByTypes indicates an expected call of CountByTypes.
func (mr *MockAccountsMockRecorder) CountByTypes(ctx, userID, _type interface{}, types ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, userID, _type}, types...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByTypes", reflect.TypeOf((*MockAccounts)(nil).CountByTypes), varargs...)
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
func (m *MockAccounts) Delete(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAccountsMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAccounts)(nil).Delete), ctx, id)
}

// Get mocks base method.
func (m *MockAccounts) Get(ctx context.Context, id int64) (domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockAccountsMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAccounts)(nil).Get), ctx, id)
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

// Update mocks base method.
func (m *MockAccounts) Update(ctx context.Context, toUpdate domain.AccountToUpdate, id int64, _type domain.AccountType) (domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, toUpdate, id, _type)
	ret0, _ := ret[0].(domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockAccountsMockRecorder) Update(ctx, toUpdate, id, _type interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAccounts)(nil).Update), ctx, toUpdate, id, _type)
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
func (m *MockTransactions) Create(ctx context.Context, toCreate domain.TransactionToCreate, categoryId, creditId, debitId *int64) (domain.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, toCreate, categoryId, creditId, debitId)
	ret0, _ := ret[0].(domain.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTransactionsMockRecorder) Create(ctx, toCreate, categoryId, creditId, debitId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTransactions)(nil).Create), ctx, toCreate, categoryId, creditId, debitId)
}

// Delete mocks base method.
func (m *MockTransactions) Delete(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTransactionsMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTransactions)(nil).Delete), ctx, id)
}

// GetOwner mocks base method.
func (m *MockTransactions) GetOwner(ctx context.Context, id int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOwner", ctx, id)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOwner indicates an expected call of GetOwner.
func (mr *MockTransactionsMockRecorder) GetOwner(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOwner", reflect.TypeOf((*MockTransactions)(nil).GetOwner), ctx, id)
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

// Stats mocks base method.
func (m *MockTransactions) Stats(ctx context.Context, filter domain.TransactionsFilter) ([]domain.TransactionStat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stats", ctx, filter)
	ret0, _ := ret[0].([]domain.TransactionStat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stats indicates an expected call of Stats.
func (mr *MockTransactionsMockRecorder) Stats(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stats", reflect.TypeOf((*MockTransactions)(nil).Stats), ctx, filter)
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

// Get mocks base method.
func (m *MockTransactionCategories) Get(ctx context.Context, id int64) (domain.TransactionCategory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(domain.TransactionCategory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockTransactionCategoriesMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTransactionCategories)(nil).Get), ctx, id)
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

// MockBalances is a mock of Balances interface.
type MockBalances struct {
	ctrl     *gomock.Controller
	recorder *MockBalancesMockRecorder
}

// MockBalancesMockRecorder is the mock recorder for MockBalances.
type MockBalancesMockRecorder struct {
	mock *MockBalances
}

// NewMockBalances creates a new mock instance.
func NewMockBalances(ctrl *gomock.Controller) *MockBalances {
	mock := &MockBalances{ctrl: ctrl}
	mock.recorder = &MockBalancesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBalances) EXPECT() *MockBalancesMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockBalances) Get(ctx context.Context, accountID int64, date time.Time) (domain.Balance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, accountID, date)
	ret0, _ := ret[0].(domain.Balance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockBalancesMockRecorder) Get(ctx, accountID, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockBalances)(nil).Get), ctx, accountID, date)
}
