# Snapshot operator
Operator for the snapshot controller and provisioner: Deploys the snapshot controller and provisioner containers into a
Kubernetes cluster.

## Build

1. Run `dep ensure`

2. Run `make`: the compiled binary is placed in `_output/bin/snahshot-operator`

## Usage
1. Create the CRD for the SnapshotController:
   ```bash
   $ kubectl create -f deploy/00-crd.yaml
   ```

2. Create the necessary RBAC objects:
   ```bash
   $ kubectl create -f deploy/01-rbac.yaml
   ```

4. Start the operator (outside of the cluster):
   ```bash
   $ KUBERNETES_CONFIG=/var/run/kubernetes/admin.kubeconfig OPERATOR_NAME="snapshot-operator" WATCH_NAMESPACE="default" _output/bin/snapshot-operator
   ```
   Make sure to use the correct value for the `KUBERNETES_CONFIG` environment variable.

5. Create the custom resource:
   ```bash
   $ kubectl create -f deploy/03-cr.yaml
   ```
