#!/bin/bash

# Copyright 2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

set -ex

if [ -z "${PRODUCT_SLUG}" ] ; then
  echo "PRODUCT_SLUG not defined"
  exit 1
fi

# Get the ID for the first container image
IMAGES=$(mkpcli container-image list --debug --debug-request-payloads --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" --output json)
IMAGE_URL=$(echo "${IMAGES}" | jq -r .dockerurlsList[0].url)
IMAGE_TAG=$(echo "${IMAGES}" | jq -r .dockerurlsList[0].imagetagsList[0].tag)
IS_IN_MKP_REGISTRY=$(echo "${IMAGES}" | jq -r .dockerurlsList[0].imagetagsList[0].isupdatedinmarketplaceregistry)
PROCESSING_ERROR=$(echo "${IMAGES}" | jq -r .dockerurlsList[0].imagetagsList[0].processingerror)

while [ "${IS_IN_MKP_REGISTRY}" == "false" ] && [ -z "${PROCESSING_ERROR}" ] ; do
  sleep 60
done

if [ "${IS_IN_MKP_REGISTRY}" == "true" ] && [ -z "${PROCESSING_ERROR}" ] ; then
  # Download the image
  mkpcli container-image download --debug --debug-request-payloads \
    --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" \
    --image-repository "${IMAGE_URL}" --tag "${IMAGE_TAG}" \
    --filename my-container-image.tar

  # Downloaded file is a real docker image
  test -f my-container-image.tar
  tar tvf my-container-image.tar manifest.json
elif [ -n "${PROCESSING_ERROR}" ] ; then
  echo "Container image is not downloadable"
  exit 1
else
  echo "Unknown status"
  exit 1
fi
