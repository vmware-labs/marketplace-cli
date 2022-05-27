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
ASSETS=$(mkpcli product list-assets --type chart --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" --output json)
NAME=$(echo "${ASSETS}" | jq -r .[0].displayname)
DOWNLOADABLE=$(echo "${ASSETS}" | jq -r .[0].downloadable)
ERROR=$(echo "${ASSETS}" | jq -r .[0].error)

if [ "${DOWNLOADABLE}" == "true" ] ; then
  mkpcli download --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" \
    --filter "${NAME}" \
    --filename my-chart.tgz \
    --accept-eula

  # Downloaded file is a real Helm chart
  test -f my-chart.tgz
  tar tvf my-chart.tgz | grep Chart.yaml

  rm -f my-chart.tgz
else
  echo "Chart is not downloadable: ${ERROR}"
  exit 1
fi
