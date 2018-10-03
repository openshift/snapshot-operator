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

	_ "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultControllerImage  = "quay.io/external_storage/snapshot-controller"
	defaultProvisionerImage = "quay.io/external_storage/snapshot-provisioner"
)

// newSnapshotController pod creates a new pod using the images in the
// SnapshotController custom resource spec
func newSnapshotControllerDeployment(cr *v1alpha1.SnapshotController) *appsv1.Deployment {
	var controllerImage string
	var provisionerImage string

	if cr.Spec.ControllerImage == "" {
		controllerImage = defaultControllerImage
	} else {
		controllerImage = cr.Spec.ControllerImage
	}
	if cr.Spec.ProvisionerImage == "" {
		provisionerImage = defaultProvisionerImage
	} else {
		provisionerImage = cr.Spec.ProvisionerImage
	}
	replicasNum := int32(1) // No other number makes sense yet

	labels := map[string]string{
		"app":              "snapshot-controller",
		"operator-managed": "true",
	}

	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "snapshot-controller",
			Namespace: cr.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicasNum,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: serviceAccountName,
					Containers: []corev1.Container{
						{
							Image: controllerImage,
							Name:  "snapshot-controller",
						},
						{
							Image: provisionerImage,
							Name:  "snapshot-provisioner",
						},
					},
				},
			},
		},
	}
	addOwnerRef(dep, ownerRefFrom(cr))

	return dep
}
