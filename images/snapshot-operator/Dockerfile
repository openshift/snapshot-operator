#oit## This file is managed by the OpenShift Image Tool: https://github.com/openshift/enterprise-images
#oit## by the OpenShift Continuous Delivery team (#aos-cd-team on IRC).
#oit## 
#oit## Any yum repos listed in this file will effectively be ignored during CD builds.
#oit## Yum repos must be enabled in the oit configuration files.
#oit## Some aspects of this file may be managed programmatically. For example, the image name, labels (version,
#oit## release, and other), and the base FROM. Changes made directly in distgit may be lost during the next
#oit## reconciliation.
#oit## 
FROM rhel7:7-released

RUN yum-config-manager --enable rhel-7-server-ose-3.7-rpms
RUN yum install -y openshift-snapshot-operator
RUN yum clean all

USER root

ENTRYPOINT ["/usr/bin/snapshot-operator"]

LABEL \
        io.k8s.description="OpenShift snapshot controller and provisioner operator." \
        com.redhat.component="snapshot-operator" \
        vendor="Red Hat" \
        description="OpenShift snapshot controller/provisioner operator" \
        License="ASL2" \
        url="https://github.com/openshift/snapshot-operator" \
        io.k8s.display-name="Snapshot operator" \
        summary="OpenShift snaphot controller/provisioner operator" \
        version="v4.0" \
        architecture="x86_64" \
        io.openshift.tags="volume,snapshot,controller,provisioner" \
        name="openshift3/origin-snapshot-operator"

