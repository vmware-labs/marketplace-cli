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

# Get the name for the first vm
FILES=$(mkpcli vm list --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" --output json)
NAME=$(echo "${FILES}" | jq -r .[0].name)
STATUS=$(echo "${FILES}" | jq -r .[0].status)

if [ "${STATUS}" == "APPROVAL_PENDING" ] || [ "${STATUS}" == "ACTIVE" ] ; then
  mkpcli download --product "${PRODUCT_SLUG}" --product-version "${PRODUCT_VERSION}" \
    --filter "${NAME}" \
    --filename my-file \
    --accept-eula

  # Downloaded virtual machine file is a real file
  test -f my-file

  rm -f my-file
elif [ "${STATUS}" == "INACTIVE" ] ; then
  echo "VM file is not downloadable"
  echo "${FILES}" | jq -r .[0].comment
  exit 1
else
  echo "Unknown status"
  exit 1
fi