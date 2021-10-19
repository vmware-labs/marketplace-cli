# Copyright 2021 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

FROM golang:1.17 as builder
ARG VERSION

COPY . /marketplace-cli/
#ENV GOPATH
ENV PATH="${PATH}:/root/go/bin"
WORKDIR /marketplace-cli/
RUN make build/mkpcli-linux-amd64

FROM photon:4.0
LABEL description="VMware Marketplace CLI"
LABEL maintainer="tanzu-isv-engineering@groups.vmware.com"

RUN yum install jq -y && \
    yum clean all

COPY --from=builder /marketplace-cli/build/mkpcli-linux-amd64 /usr/local/bin/mkpcli
ENTRYPOINT ["/usr/local/bin/mkpcli"]
