/*
   Copyright 2018 Red Hat, Inc.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package operator

import (
	"github.com/openshift/snapshot-operator/pkg/apis/snapshotoperator/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	serviceAccountName = "snapshot-controller-runner"
	clusterRoleName    = "snapshot-controller-role"
	roleBindingName    = "snapshot-controller"
)

func newServiceAccount(cr *v1alpha1.SnapshotController) *corev1.ServiceAccount {
	acc := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: cr.Namespace,
		},
	}
	addOwnerRef(acc, ownerRefFrom(cr))

	return acc
}

func newClusterRole(cr *v1alpha1.SnapshotController) *rbacv1.ClusterRole {
	role := &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRole",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterRoleName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs: []string{
					"get",
					"list",
					"watch",
					"create",
					"delete",
				},
				APIGroups: []string{""},
				Resources: []string{"persistentvolumes"},
			},
			{
				Verbs: []string{
					"get",
					"list",
					"watch",
					"update",
				},
				APIGroups: []string{""},
				Resources: []string{"persistentvolumeclaims"},
			},
			{
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
				APIGroups: []string{"storage.k8s.io"},
				Resources: []string{"storageclasses"},
			},
			{
				Verbs: []string{
					"list",
					"watch",
					"create",
					"update",
					"patch",
				},
				APIGroups: []string{""},
				Resources: []string{"events"},
			},
			{
				Verbs: []string{
					"create",
					"list",
					"watch",
					"delete",
				},
				APIGroups: []string{"apiextensions.k8s.io"},
				Resources: []string{"customresourcedefinitions"},
			},
			{
				Verbs: []string{
					"get",
					"list",
					"watch",
					"create",
					"update",
					"patch",
					"delete",
				},
				APIGroups: []string{"volumesnapshot.external-storage.k8s.io"},
				Resources: []string{
					"volumesnapshots",
					"volumesnapshotdatas",
				},
			},
		},
	}
	addOwnerRef(role, ownerRefFrom(cr))

	return role
}

func newClusterRoleBinding(cr *v1alpha1.SnapshotController) *rbacv1.ClusterRoleBinding {

	binding := &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRoleBinding",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: roleBindingName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccountName,
				Namespace: cr.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind: "ClusterRole",
			Name: clusterRoleName,
		},
	}

	addOwnerRef(binding, ownerRefFrom(cr))

	return binding
}
