#!/bin/bash

# Copyright 2023 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

set -ex

if [ -z "${ASSET_TYPE}" ]; then
  echo "ASSET_TYPE not defined. Should be one of: chart, image, metafile, other, vm"
  exit 1
fi

if [ -z "${PRODUCT_SLUG}" ] ; then
  echo "PRODUCT_SLUG not defined"
  exit 1
fi

if [ -z "${PRODUCT_VERSION}" ] ; then
  echo "PRODUCT_VERSION not defined, using latest version" >&2
fi

filename=""
if [ "${ASSET_TYPE}" == "chart" ]; then
  filename="my-chart.tgz"
elif [ "${ASSET_TYPE}" == "image" ]; then
  filename="my-container-image.tar"
elif [ "${ASSET_TYPE}" == "metafile" ]; then
  filename="my-metafile"
elif [ "${ASSET_TYPE}" == "other" ]; then
  filename="my-addon.vlcp"
elif [ "${ASSET_TYPE}" == "vm" ]; then
  filename="my-vm-image"
fi

# Get the ID for the asset of type
ASSETS=$(mkpcli product list-assets --type "${ASSET_TYPE}" --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" --output json)
NAME=$(echo "${ASSETS}" | jq -r .[0].displayname)
DOWNLOADABLE=$(echo "${ASSETS}" | jq -r .[0].downloadable)
ERROR=$(echo "${ASSETS}" | jq -r .[0].error)

if [ "${DOWNLOADABLE}" == "true" ] ; then
  mkpcli download --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" \
    --filter "${NAME}" --filename "${filename}" --accept-eula

  test -f "${filename}"
#  if [ "${ASSET_TYPE}" == "addon" ]; then
  if [ "${ASSET_TYPE}" == "chart" ]; then
    tar tvf "${filename}" | grep Chart.yaml
  elif [ "${ASSET_TYPE}" == "image" ]; then
    tar tvf "${filename}" manifest.json
#  elif [ "${ASSET_TYPE}" == "metafile" ]; then
#  elif [ "${ASSET_TYPE}" == "vm" ]; then
  fi

  rm -f "${filename}"
else
  echo "Asset is not downloadable: ${ERROR}"
  exit 1
fi
