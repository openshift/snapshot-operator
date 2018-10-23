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
	"context"
	"encoding/json"

	"github.com/golang/glog"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/snapshot-operator/pkg/apis/snapshotoperator/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
)

func NewHandler() sdk.Handler {
	platform, err := getPlatform()
	if err != nil {
		// Migh not be running in OpenShift so jut log an error
		glog.Errorf("Error querying the cluster platform: %v", err)
	}
	handler := &Handler{
		cfg: Config{
			cloudProvider: "",
			awsSecret:     "",
		},
	}

	if platform != nil {
		// Some platform names correspond to cloud providers
		if platform.Name() == "aws" || platform.Name() == "openstack" || platform.Name() == "gce" {
			handler.cfg.cloudProvider = platform.Name()
		}
	}
	// TODO: Find cluster default AWS credentials if not overriden in the CR

	return handler
}

type Handler struct {
	cfg Config
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch cr := event.Object.(type) {
	case *v1alpha1.SnapshotController:
		// Ignore the delete event since the garbage collector will clean up all secondary resources for the CR
		// All secondary resources must have the CR set as their OwnerReference for this to be the case
		if event.Deleted {
			glog.V(4).Info("SnapshotController CR deleted")
			return nil
		}

		// Figure out where to look for the AWS secrets if needed
		if h.cfg.cloudProvider == "aws" {
			if cr.Spec.AWSSecretName != "" {
				h.cfg.awsSecret = cr.Spec.AWSSecretName
			}
			if cr.Spec.AWSSecretNamespace != "" {
				h.cfg.awsSecretNamespace = cr.Spec.AWSSecretNamespace
			} else {
				h.cfg.awsSecretNamespace = cr.Namespace
			}
		}

		// New SnapshotController CR: create RBAC and the deployment
		err := createObjectIfNotExist(newServiceAccount(cr))
		if err != nil {
			return err
		}
		err = createObjectIfNotExist(newClusterRole(cr))
		if err != nil {
			return err
		}
		err = createObjectIfNotExist(newClusterRoleBinding(cr))
		if err != nil {
			return err
		}
		if cr.Spec.StorageClassName != "" {
			err = createObjectIfNotExist(newStorageClass(cr))
			if err != nil {
				return err
			}
		}
		err = createObjectIfNotExist(newSnapshotControllerDeployment(cr, h.cfg))
		if err != nil {
			return err
		}

	}
	return nil
}

func getPlatform() (*types.Platform, error) {
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster-config-v1",
			Namespace: "kube-system",
		},
	}
	err := sdk.Get(cm)
	if err != nil {
		return nil, err
	}

	data, err := utilyaml.ToJSON([]byte(cm.Data["install-config"]))
	if err != nil {
		return nil, err
	}

	config := &types.InstallConfig{}
	json.Unmarshal(data, &config)

	return &config.Platform, nil
}
