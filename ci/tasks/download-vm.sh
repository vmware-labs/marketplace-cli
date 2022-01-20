#!/bin/bash

# Copyright 2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

set -ex

if [ -z "${PRODUCT_SLUG}" ] ; then
  echo "PRODUCT_SLUG not defined"
  exit 1
fi

# Get the name for the first vm
FILES=$(mkpcli vm list --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" --output json)
FILE_NAME=$(echo "${FILES}" | jq -r .[0].filename)
STATUS=$(echo "${FILES}" | jq -r .[0].status)

if [ "${STATUS}" == "APPROVAL_PENDING" ] || [ "${STATUS}" == "ACTIVE" ] ; then
  # Download the file
  mkpcli vm download --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" \
    --filter "${FILE_NAME}" \
    --filename my-file

  # Downloaded virtual machine file is a real file
  test -f my-file
elif [ "${STATUS}" == "INACTIVE" ] ; then
  echo "VM file is not downloadable"
  echo "${FILES}" | jq -r .[0].comment
  exit 1
else
  echo "Unknown status"
  exit 1
fi
