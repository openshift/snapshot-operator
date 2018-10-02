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

	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newStorageClass(cr *v1alpha1.SnapshotController) *storagev1.StorageClass {
	sc := &storagev1.StorageClass{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "storage.k8s.io/v1",
			Kind:       "StorageClass",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: cr.Spec.StorageClassName,
		},
		Provisioner: "volumesnapshot.external-storage.k8s.io/snapshot-promoter",
	}

	addOwnerRef(sc, ownerRefFrom(cr))
	return sc
}
