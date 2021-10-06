#!/bin/bash

# Copyright 2021 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

set -ex

if [ -z "${PRODUCT_SLUG}" ] ; then
  echo "PRODUCT_SLUG not defined"
  exit 1
fi

# Get the ID for the first chart
CHARTS=$(mkpcli chart list --debug --debug-request-payloads --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" --output json)
CHART_ID=$(echo "${CHARTS}" | jq -r .[0].id)
IS_IN_MKP_REGISTRY=$(echo "${IMAGES}" | jq -r .[0].isupdatedinmarketplaceregistry)
PROCESSING_ERROR=$(echo "${CHARTS}" | jq -r .[0].processingerror)

while [ "${IS_IN_MKP_REGISTRY}" == "false" ] && [ -z "${PROCESSING_ERROR}" ] ; do
  sleep 60
done

if [ "${IS_IN_MKP_REGISTRY}" == "true" ] ; then
  # Download the chart
  mkpcli chart download --debug --debug-request-payloads \
    --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" \
    --chart-id "${CHART_ID}" \
    --filename my-chart.tgz

  # Downloaded file is a real Helm chart
  test -f my-chart.tgz
  tar tvf my-chart.tgz | grep Chart.yaml
elif [ "${IS_IN_MKP_REGISTRY}" == "false" ] && [ -n "${PROCESSING_ERROR}" ] ; then
  echo "Chart is not downloadable"
  exit 1
else
  echo "Unknown status"
  exit 1
fi
