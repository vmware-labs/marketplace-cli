#!/usr/bin/env bash
function TestDownloadChart() {
    echo "Download Chart test"
    PRODUCT_SLUG=$(jq -r .tests.downloadChart.slug ./external-test-inputs.json)
    PRODUCT_VERSION=$(jq -r .tests.downloadChart.version ./external-test-inputs.json)
    export PRODUCT_SLUG PRODUCT_VERSION
    ./download-chart.sh
}

function TestDownloadContainerImage() {
  echo "Download Container Image test"
  PRODUCT_SLUG=$(jq -r .tests.downloadContainerImage.slug ./external-test-inputs.json)
  PRODUCT_VERSION=$(jq -r .tests.downloadContainerImage.version ./external-test-inputs.json)
  export PRODUCT_SLUG PRODUCT_VERSION
  ./download-container-image.sh
}

function TestDownloadISO() {
  echo "Download ISO test"
  PRODUCT_SLUG=$(jq -r .tests.downloadISO.slug ./external-test-inputs.json)
  PRODUCT_VERSION=$(jq -r .tests.downloadISO.version ./external-test-inputs.json)
  export PRODUCT_SLUG PRODUCT_VERSION
  ./download-vm.sh
}

function TestDownloadOVA() {
  echo "Download OVA test"
  PRODUCT_SLUG=$(jq -r .tests.downloadOVA.slug ./external-test-inputs.json)
  PRODUCT_VERSION=$(jq -r .tests.downloadOVA.version ./external-test-inputs.json)
  export PRODUCT_SLUG PRODUCT_VERSION
  ./download-vm.sh
}

TestDownloadChart
TestDownloadContainerImage
TestDownloadISO
TestDownloadOVA