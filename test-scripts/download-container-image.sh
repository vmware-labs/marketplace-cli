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
IMAGES=$(mkpcli container-image list --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" --output json)
IMAGE_URL=$(echo "${IMAGES}" | jq -r .[0].dockerurlsList[0].url)
IMAGE_TAG=$(echo "${IMAGES}" | jq -r .[0].dockerurlsList[0].imagetagsList[0].tag)
IS_IN_MKP_REGISTRY=$(echo "${IMAGES}" | jq -r .[0].dockerurlsList[0].imagetagsList[0].isupdatedinmarketplaceregistry)
PROCESSING_ERROR=$(echo "${IMAGES}" | jq -r .[0].dockerurlsList[0].imagetagsList[0].processingerror)

if [ "${IS_IN_MKP_REGISTRY}" == "true" ] ; then
  # Download the image
  mkpcli download --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" \
    --filter "${IMAGE_URL}:${IMAGE_TAG}" \
    --filename my-container-image.tar \
    --accept-eula

  # Downloaded file is a real docker image
  test -f my-container-image.tar
  tar tvf my-container-image.tar manifest.json

  rm -f my-container-image.tar
elif [ -n "${PROCESSING_ERROR}" ] ; then
  echo "Container image is not downloadable"
  exit 1
else
  echo "Unknown status"
  exit 1
fi
