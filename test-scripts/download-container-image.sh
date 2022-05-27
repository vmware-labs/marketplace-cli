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

# Get the ID for the first container image
ASSETS=$(mkpcli product list-assets --type image --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" --output json)
Name=$(echo "${ASSETS}" | jq -r .[0].displayname)
DOWNLOADABLE=$(echo "${ASSETS}" | jq -r .[0].downloadable)
ERROR=$(echo "${ASSETS}" | jq -r .[0].error)

if [ "${DOWNLOADABLE}" == "true" ] ; then
  mkpcli download --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" \
    --filter "${Name}" \
    --filename my-container-image.tar \
    --accept-eula

  # Downloaded file is a real container image
  test -f my-container-image.tar
  tar tvf my-container-image.tar manifest.json

  rm -f my-container-image.tar
else
  echo "Container image is not downloadable: ${ERROR}"
  exit 1
fi
