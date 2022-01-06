# Updating Products

A primary focus of the Marketplace CLI is for updating a product in the VMware Marketplace.

While the CLI is not meant for creating new products, it is excellent at creating new versions and attaching binaries to the product.

## What can be updated

Currently, the Marketplace CLI can be used to:

* Add version numbers
* Attach container images
* Attach Helm charts
* Attach virtual machine files (ISOs & OVAs)

## Versions

Adding a version number to a product:

```bash
$ mkpcli product add-version --product hyperspace-database --version 1.0.1
Versions for Hyperspace Database:
  NUMBER  STATUS
  1.0.1   PENDING
  1.0.0   PENDING
  0.3.0   PENDING
```

## Container Images

Attaching a container image to a product:

```bash
mkpcli container-image attach --product hyperspace-database --version 1.0.1 --image-repository astrowidgets/hyperspacedb --tag 1.0.1 --tag-type FIXED --deployment-instructions 'docker run astrowidgets/hyperspacedb:1.0.1'
```

NOTE: The image will be pulled by the Marketplace servers and stored locally, so it must be publicly reachable.

## Helm Charts

Attaching a container image to a product:

### Upload a local chart

```bash
mkpcli chart attach --product hyperspace-database-chart --product-version 1.0.1 --chart charts/hyperspace-db-1.0.1.tgz --readme 'helm install it'
```

### Attaching a remote chart

```bash
mkpcli chart attach --product hyperspace-database-chart --product-version 1.0.1 --chart https://astro-widgets.example.com/charts/hyperspace-db-1.0.1.tgz --readme 'helm install it'
```

NOTE: The chart will be pulled by the Marketplace servers and stored locally, so it must be publicly reachable.

## VM Images (e.g. ISOs, OVAs)

Uploading and attaching a VM image to a product:

```bash
mkpcli vm attach --product hyperspace-database-vm --product-version 1.0.1 --file vm/hyperspace-db-1.0.1-1526e30ba.iso
```
