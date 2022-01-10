# Concourse Example

Running the Marketplace CLI inside [Concourse](https://concourse-ci.org/) is simple.

You will need the version of the product, which is typically handled with the [semver resource](https://github.com/concourse/semver-resource), and the actual asset to attach (container image, Helm chart, VM ISO or OVA).

Your product slug and the [CSP API token](Authentication.md) can be passed via parameters.

## Example pipeline
```yaml
resource:
  - name: hyperspace-database
    type: helm-chart

  - name: version
    type: semver

  - name: mkpcli
    type: registry-image
    source:
      repository: projects.registry.vmware.com/tanzu_isv_engineering/mkpcli    

jobs:
  - name: update-marketplace-product
    plan:
      - in_parallel:
          - get: hyperspace-db-chart
          - get: version
          - get: mkpcli
      - task: add-version
        image: mkpcli
        config:
          params:
            CSP_API_TOKEN: ((csp.api_token))
            PRODUCT_SLUG: hyperspace-database
          platform: linux
          inputs:
            - name: version
          run:
              path: bash
              args:
                - -xc
                - |
                  export VERSION=$(cat version/version)
            
                  mkpcli product list-versions --product "${PRODUCT_SLUG}" | grep "${VERSION}"
                  if [[ $? -ne 0 ]] ; then
                    set -e
                    mkpcli product add-version --product "${PRODUCT_SLUG}" --product-version "${VERSION}"
                    mkpcli product list-versions --product "${PRODUCT_SLUG}" | grep "${VERSION}"
                  fi
      - task: add-chart
        image: mkpcli
        config:
          params:
            CSP_API_TOKEN: ((csp.api_token))
            PRODUCT_SLUG: hyperspace-database
          platform: linux
          inputs:
            - name: hyperspace-db-chart
            - name: version
          run:
              path: bash
              args:
                - -exc
                - |
                  VERSION=$(cat version/version)

                  mkpcli chart list --product "${PRODUCT_SLUG}" --product-version "${VERSION}"
                  mkpcli chart attach --product "${PRODUCT_SLUG}" --product-version "${VERSION}" --chart chart/*.tgz --readme "helm install it"            
                  mkpcli chart list --product "${PRODUCT_SLUG}" --product-version "${VERSION}"
```