# Marketplace CLI

`mkpcli` enables a command-line interface to the [VMware Marketplace](http://marketplace.cloud.vmware.com/) for consumers and publishes.

To install, grab the latest prebuilt binary from the [Releases](https://github.com/vmware-labs/marketplace-cli/releases) page, or [build from source](#building).

Features:
* Get details about a product
* Manage products in your org
  * Add versions
  * Attach container images
  * Attach Helm charts
  * Attach virtual machine files (ISOs & OVAs)
* Download assets from a product

## Example
```bash
$ export CSP_API_TOKEN=...
$ mkpcli product add-version --product hyperspace-database --version 1.0.1
$ mkpcli chart attach --product hyperspace-database --version 1.0.1 --chart ./hyperspace-database-1.0.1.tgz
```

For more information, see [Updating Products](docs/UpdatingProducts.md)


## Authentication

`mkpcli` requires an API Token from [VMware Cloud Services](https://console.cloud.vmware.com/csp/gateway/portal/#/user/tokens). See [this doc](./docs/Authentication.md) for more information.

The token can be set via CLI flag (i.e. `--csp-api-token`) or environment variable (i.e. `CSP_API_TOKEN`).

For more information, see [Authentication](docs/Authentication.md)

## Building

Building from source is simple with our Makefile:

```bash
$ make build
...
go build -o build/mkpcli -ldflags "-X github.com/vmware-labs/marketplace-cli/v2/cmd.version=dev" ./main.go
$ file build/mkpcli 
build/mkpcli: Mach-O 64-bit executable x86_64
$ ./build/mkpcli 
mkpcli is a CLI interface for the VMware Marketplace,
enabling users to view, get, and manage their Marketplace entries.
...
```

## Developing

If you would like to build and contribute to this project, please fork and make pull requests.

If you are internal to VMware, and you would like to run commands against the [Marketplace staging service](https://stg.market.csp.vmware.com/), set this environment variable:
```
export MARKETPLACE_ENV=staging
```

Please see our [Code of Conduct](CODE-OF-CONDUCT.md) and [Contributors guide](CONTRIBUTING.md).

