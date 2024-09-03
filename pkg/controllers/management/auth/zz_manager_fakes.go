// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go
//
// Generated by this command:
//
//	mockgen -source=manager.go -destination=zz_manager_fakes.go -package=auth
//

// Package auth is a generated GoMock package.
package auth

import (
	reflect "reflect"

	v3 "github.com/rancher/rancher/pkg/generated/norman/management.cattle.io/v3"
	gomock "go.uber.org/mock/gomock"
	v1 "k8s.io/api/rbac/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// MockmanagerInterface is a mock of managerInterface interface.
type MockmanagerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockmanagerInterfaceMockRecorder
}

// MockmanagerInterfaceMockRecorder is the mock recorder for MockmanagerInterface.
type MockmanagerInterfaceMockRecorder struct {
	mock *MockmanagerInterface
}

// NewMockmanagerInterface creates a new mock instance.
func NewMockmanagerInterface(ctrl *gomock.Controller) *MockmanagerInterface {
	mock := &MockmanagerInterface{ctrl: ctrl}
	mock.recorder = &MockmanagerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmanagerInterface) EXPECT() *MockmanagerInterfaceMockRecorder {
	return m.recorder
}

// checkReferencedRoles mocks base method.
func (m *MockmanagerInterface) checkReferencedRoles(arg0, arg1 string, arg2 int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "checkReferencedRoles", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// checkReferencedRoles indicates an expected call of checkReferencedRoles.
func (mr *MockmanagerInterfaceMockRecorder) checkReferencedRoles(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "checkReferencedRoles", reflect.TypeOf((*MockmanagerInterface)(nil).checkReferencedRoles), arg0, arg1, arg2)
}

// ensureClusterMembershipBinding mocks base method.
func (m *MockmanagerInterface) ensureClusterMembershipBinding(arg0, arg1 string, arg2 *v3.Cluster, arg3 bool, arg4 v1.Subject) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ensureClusterMembershipBinding", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// ensureClusterMembershipBinding indicates an expected call of ensureClusterMembershipBinding.
func (mr *MockmanagerInterfaceMockRecorder) ensureClusterMembershipBinding(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ensureClusterMembershipBinding", reflect.TypeOf((*MockmanagerInterface)(nil).ensureClusterMembershipBinding), arg0, arg1, arg2, arg3, arg4)
}

// ensureProjectMembershipBinding mocks base method.
func (m *MockmanagerInterface) ensureProjectMembershipBinding(arg0, arg1, arg2 string, arg3 *v3.Project, arg4 bool, arg5 v1.Subject) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ensureProjectMembershipBinding", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// ensureProjectMembershipBinding indicates an expected call of ensureProjectMembershipBinding.
func (mr *MockmanagerInterfaceMockRecorder) ensureProjectMembershipBinding(arg0, arg1, arg2, arg3, arg4, arg5 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ensureProjectMembershipBinding", reflect.TypeOf((*MockmanagerInterface)(nil).ensureProjectMembershipBinding), arg0, arg1, arg2, arg3, arg4, arg5)
}

// grantManagementClusterScopedPrivilegesInProjectNamespace mocks base method.
func (m *MockmanagerInterface) grantManagementClusterScopedPrivilegesInProjectNamespace(arg0, arg1 string, arg2 map[string]string, arg3 v1.Subject, arg4 *v3.ClusterRoleTemplateBinding) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "grantManagementClusterScopedPrivilegesInProjectNamespace", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// grantManagementClusterScopedPrivilegesInProjectNamespace indicates an expected call of grantManagementClusterScopedPrivilegesInProjectNamespace.
func (mr *MockmanagerInterfaceMockRecorder) grantManagementClusterScopedPrivilegesInProjectNamespace(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "grantManagementClusterScopedPrivilegesInProjectNamespace", reflect.TypeOf((*MockmanagerInterface)(nil).grantManagementClusterScopedPrivilegesInProjectNamespace), arg0, arg1, arg2, arg3, arg4)
}

// grantManagementPlanePrivileges mocks base method.
func (m *MockmanagerInterface) grantManagementPlanePrivileges(arg0 string, arg1 map[string]string, arg2 v1.Subject, arg3 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "grantManagementPlanePrivileges", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// grantManagementPlanePrivileges indicates an expected call of grantManagementPlanePrivileges.
func (mr *MockmanagerInterfaceMockRecorder) grantManagementPlanePrivileges(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "grantManagementPlanePrivileges", reflect.TypeOf((*MockmanagerInterface)(nil).grantManagementPlanePrivileges), arg0, arg1, arg2, arg3)
}

// grantManagementProjectScopedPrivilegesInClusterNamespace mocks base method.
func (m *MockmanagerInterface) grantManagementProjectScopedPrivilegesInClusterNamespace(arg0, arg1 string, arg2 map[string]string, arg3 v1.Subject, arg4 *v3.ProjectRoleTemplateBinding) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "grantManagementProjectScopedPrivilegesInClusterNamespace", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// grantManagementProjectScopedPrivilegesInClusterNamespace indicates an expected call of grantManagementProjectScopedPrivilegesInClusterNamespace.
func (mr *MockmanagerInterfaceMockRecorder) grantManagementProjectScopedPrivilegesInClusterNamespace(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "grantManagementProjectScopedPrivilegesInClusterNamespace", reflect.TypeOf((*MockmanagerInterface)(nil).grantManagementProjectScopedPrivilegesInClusterNamespace), arg0, arg1, arg2, arg3, arg4)
}

// reconcileClusterMembershipBindingForDelete mocks base method.
func (m *MockmanagerInterface) reconcileClusterMembershipBindingForDelete(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "reconcileClusterMembershipBindingForDelete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// reconcileClusterMembershipBindingForDelete indicates an expected call of reconcileClusterMembershipBindingForDelete.
func (mr *MockmanagerInterfaceMockRecorder) reconcileClusterMembershipBindingForDelete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "reconcileClusterMembershipBindingForDelete", reflect.TypeOf((*MockmanagerInterface)(nil).reconcileClusterMembershipBindingForDelete), arg0, arg1)
}

// reconcileProjectMembershipBindingForDelete mocks base method.
func (m *MockmanagerInterface) reconcileProjectMembershipBindingForDelete(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "reconcileProjectMembershipBindingForDelete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// reconcileProjectMembershipBindingForDelete indicates an expected call of reconcileProjectMembershipBindingForDelete.
func (mr *MockmanagerInterfaceMockRecorder) reconcileProjectMembershipBindingForDelete(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "reconcileProjectMembershipBindingForDelete", reflect.TypeOf((*MockmanagerInterface)(nil).reconcileProjectMembershipBindingForDelete), arg0, arg1, arg2)
}

// removeAuthV2Permissions mocks base method.
func (m *MockmanagerInterface) removeAuthV2Permissions(arg0 string, arg1 runtime.Object) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "removeAuthV2Permissions", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// removeAuthV2Permissions indicates an expected call of removeAuthV2Permissions.
func (mr *MockmanagerInterfaceMockRecorder) removeAuthV2Permissions(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "removeAuthV2Permissions", reflect.TypeOf((*MockmanagerInterface)(nil).removeAuthV2Permissions), arg0, arg1)
}