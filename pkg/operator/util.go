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

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// objects can have more than one ownerRef, potentially
func addOwnerRef(obj metav1.Object, ownerRef *metav1.OwnerReference) {
	if obj != nil {
		if ownerRef != nil {
			obj.SetOwnerReferences(append(obj.GetOwnerReferences(), *ownerRef))
		}
	}
}

func ownerRefFrom(cr *v1alpha1.SnapshotController) *metav1.OwnerReference {
	if cr != nil {
		truthy := true
		return &metav1.OwnerReference{
			APIVersion: cr.APIVersion,
			Kind:       cr.Kind,
			Name:       cr.Name,
			UID:        cr.UID,
			Controller: &truthy,
		}
	}
	return nil
}

// helper function to create an API object: wiil just log a warning if
// the object already exists
func createObjectIfNotExist(o sdk.Object) error {
	err := sdk.Create(o)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			logrus.Warnf("failed to create snapshot controller/provisioner deployment (will try to continue): %v", err)

		} else {
			logrus.Errorf("failed to create snapshot controller/provisioner deployment: %v", err)
			return err
		}
	}

	return nil
}