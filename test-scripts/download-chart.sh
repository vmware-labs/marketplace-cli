#!/bin/bash

# Copyright 2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

set -ex

if [ -z "${PRODUCT_SLUG}" ] ; then
  echo "PRODUCT_SLUG not defined"
  exit 1
fi

if [ -z "${PRODUCT_VERSION}" ] ; then
  echo "PRODUCT_VERSION not defined, using latest version" >&2
fi

# Get the name for the first chart
CHARTS=$(mkpcli chart list --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" --output json)
CHART_NAME=$(echo "${CHARTS}" | jq -r .[0].helmtarurl)
IS_IN_MKP_REGISTRY=$(echo "${CHARTS}" | jq -r .[0].isupdatedinmarketplaceregistry)
PROCESSING_ERROR=$(echo "${CHARTS}" | jq -r .[0].processingerror)

while [ "${IS_IN_MKP_REGISTRY}" == "false" ] && [ -z "${PROCESSING_ERROR}" ] ; do
  sleep 60
done

if [ "${IS_IN_MKP_REGISTRY}" == "true" ] ; then
  mkpcli download --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" \
    --filter "${CHART_NAME}" \
    --filename my-chart.tgz \
    --accept-eula

  # Downloaded file is a real Helm chart
  test -f my-chart.tgz
  tar tvf my-chart.tgz | grep Chart.yaml

  rm -f my-chart.tgz
elif [ "${IS_IN_MKP_REGISTRY}" == "false" ] && [ -n "${PROCESSING_ERROR}" ] ; then
  echo "Chart is not downloadable"
  exit 1
else
  echo "Unknown status"
  exit 1
fi
