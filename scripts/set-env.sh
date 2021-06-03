#!/bin/bash

# Copyright 2021 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

echo CSP_API_TOKEN
CSP_API_TOKEN=$(vault read /runway_concourse/tanzu-isv-engineering/marketplace_api_token -format=json | jq -r .data.value)

echo MARKETPLACE_HOST
MARKETPLACE_HOST="gtwstg.market.csp.vmware.com"
if [ "$1" == "production" ] ; then
  MARKETPLACE_HOST="gtw.marketplace.cloud.vmware.com"
fi

export CSP_API_TOKEN \
    MARKETPLACE_HOST
