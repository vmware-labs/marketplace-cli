FROM harbor-repo.vmware.com/dockerhub-proxy-cache/library/golang:1.16
RUN apt-get update && \
    apt-get install -y zip && \
    rm -rf /var/lib/apt/lists/*

ARG PROTOC_VERSION=3.14.0
RUN PROTOC_ZIP=protoc-${PROTOC_VERSION}-linux-x86_64.zip && \
    curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/$PROTOC_ZIP && \
    unzip -o $PROTOC_ZIP -d /usr/local bin/protoc && \
    unzip -o $PROTOC_ZIP -d /usr/local 'include/*' && \
    rm -f $PROTOC_ZIP

RUN go get google.golang.org/protobuf/cmd/protoc-gen-go && \
    go get github.com/favadi/protoc-go-inject-tag && \
    go get golang.org/x/tools/cmd/goimports
