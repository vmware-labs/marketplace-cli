#!/bin/bash

echo CSP_API_TOKEN
CSP_API_TOKEN=$(vault read /runway_concourse/tanzu-isv-engineering/marketplace_api_token -format=json | jq -r .data.value)
GO111MODULE=on
GOPRIVATE=gitlab.eng.vmware.com

echo MARKETPLACE_HOST
MARKETPLACE_HOST="gtwstg.market.csp.vmware.com"
if [ "$1" == "production" ] ; then
  MARKETPLACE_HOST="gtw.marketplace.cloud.vmware.com"
fi

export CSP_API_TOKEN \
    GO111MODULE \
    GOPRIVATE \
    MARKETPLACE_HOST
