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

	"github.com/golang/glog"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
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
func newSnapshotControllerDeployment(cr *v1alpha1.SnapshotController, cfg Config) *appsv1.Deployment {
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
	// Add cloudprovider commandline flag
	if cfg.cloudProvider != "" {
		for _, container := range dep.Spec.Template.Spec.Containers {
			container.Args = []string{"-cloudprovider", cfg.cloudProvider}
		}
	}
	// Need a secret with credentials for AWS
	if cfg.cloudProvider == "aws" {
		if cfg.awsSecret != "" {
			if cfg.awsSecretNamespace != cr.Namespace {
				// Copy the AWS secret to the CR namespace so it's acessible by the pods
				_, err := copySecret(cfg.awsSecret, cfg.awsSecretNamespace, cr.Namespace)
				if err != nil {
					glog.Warning("Error copying AWS secret to deployment namespace: %v", err)
					return nil
				}
			}
			for _, container := range dep.Spec.Template.Spec.Containers {
				container.Env = []corev1.EnvVar{
					{
						Name: "AWS_ACCESS_KEY_ID",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{Name: cfg.awsSecret},
								Key:                  "access-key-id",
							},
						},
					},
					{
						Name: "AWS_SECRET_ACCESS_KEY",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{Name: cfg.awsSecret},
								Key:                  "secret-access-key",
							},
						},
					},
				}
			}
		} else {
			glog.Warning("AWS cloud provider requested but AWS credentials not found")
		}

	}
	addOwnerRef(dep, ownerRefFrom(cr))

	return dep
}

func copySecret(secretName string, fromNamespace string, toNamespace string) (*corev1.Secret, error) {
	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: fromNamespace,
		},
	}
	err := sdk.Get(secret)
	if err != nil {
		return nil, err
	}
	secret.Namespace = toNamespace
	err = sdk.Create(secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}
