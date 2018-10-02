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
	"github.com/openshift/snapshot-operator/pkg/apis/snapshotoperator/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	logrus.Info("Called event handler!")
	switch cr := event.Object.(type) {
	case *v1alpha1.SnapshotController:
		// Ignore the delete event since the garbage collector will clean up all secondary resources for the CR
		// All secondary resources must have the CR set as their OwnerReference for this to be the case
		if event.Deleted {
			return nil
		}

		// New SnapshotController CR: create RBAC and the deployment
		logrus.Info("Creating service account")
		//err := createObjectIfNotExist(newServiceAccount(cr))
		err := sdk.Create(newServiceAccount(cr))
		if err != nil {
			return err
		}
		logrus.Info("Creating ClusterRole")
		err = createObjectIfNotExist(newClusterRole(cr))
		if err != nil {
			return err
		}
		logrus.Info("Creating ClusterRoleBinding")
		err = createObjectIfNotExist(newClusterRoleBinding(cr))
		if err != nil {
			return err
		}
		logrus.Info("Creating StorageClass")
		if cr.Spec.StorageClassName != "" {
			err = createObjectIfNotExist(newStorageClass(cr))
			if err != nil {
				return err
			}
		}
		logrus.Info("Creating Deployment")
		err = createObjectIfNotExist(newSnapshotControllerDeployment(cr))
		if err != nil {
			return err
		}

	}
	return nil
}
