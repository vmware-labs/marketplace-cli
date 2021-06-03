# Copyright 2021 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

FROM harbor-repo.vmware.com/dockerhub-proxy-cache/library/ubuntu
LABEL description="VMware Marketplace CLI"
LABEL maintainer="tanzu-isv-engineering@groups.vmware.com"

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY build/mkpcli-linux /usr/local/bin/mkpcli
ENTRYPOINT ["/usr/local/bin/mkpcli"]
