# Copyright 2021 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

FROM harbor-repo.vmware.com/dockerhub-proxy-cache/library/photon:4.0
LABEL description="VMware Marketplace CLI"
LABEL maintainer="tanzu-isv-engineering@groups.vmware.com"

RUN yum install jq -y && \
    yum clean all

COPY build/mkpcli-linux /usr/local/bin/mkpcli
ENTRYPOINT ["/usr/local/bin/mkpcli"]
