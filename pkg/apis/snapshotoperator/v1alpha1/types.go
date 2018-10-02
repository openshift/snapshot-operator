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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SnapshotController is the main CR being watched by the operator
type SnapshotControllerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []SnapshotController `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SnapshotController is the main CR being watched by the operator
type SnapshotController struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   SnapshotControllerSpec   `json:"spec"`
	Status SnapshotControllerStatus `json:"status,omitempty"`
}

// SnapshotControllerSpec defines the configuration of the volume
// snapshot controller and provisioner
type SnapshotControllerSpec struct {
	// Default controller image location override
	// Optional, default empty
	ControllerImage string `json: "controllerImage",omitempty`

	// Default provisioner image location override
	// Optional, default empty
	ProvisionerImage string `json: "provisionerImage",omitempty`

	// Name of storage class to create. No storage class will be
	// created when the name is empty. It makes no sense for this
	// StorageClass to be a default one, this can’t be configured.
	// Optional, no default.
	StorageClassName string `json: "storageClassName",omitempty`

	// Group that would be allowed to use volume snapshots
	// Optional, default empty
	Group string `json: "group",omitempty`
}

type SnapshotControllerStatus struct {
	//// TODO: Figure out what belongs here...
}
