#!/bin/bash

# Copyright 2021 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

set -ex

if [ -z "${PRODUCT_SLUG}" ] ; then
  echo "PRODUCT_SLUG not defined"
  exit 1
fi

# Get the ID for the first chart
FILES=$(mkpcli ova list --debug --debug-request-payloads --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" --output json)
FILE_ID=$(echo "${FILES}" | jq -r .[0].fileid)
STATUS=$(echo "${FILES}" | jq -r .[0].status)

#while [ "${VALIDATION_STATUS}" == "pending" ] && [ -z "${PROCESSING_ERROR}" ] ; do
#  sleep 60
#done

if [ "${STATUS}" == "APPROVAL_PENDING" ] || [ "${STATUS}" == "ACTIVE" ] ; then
  # Download the file
  mkpcli ova download --debug --debug-request-payloads \
    --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" \
    --file-id "${FILE_ID}" \
    --filename my-file.ova

  # Downloaded file is a real OVA
  test -f my-file.ova
elif [ "${STATUS}" == "INACTIVE" ] ; then
  echo "Chart is not downloadable"
  echo "${FILES}" | jq -r .[0].comment
  exit 1
else
  echo "Unknown status"
  exit 1
fi
