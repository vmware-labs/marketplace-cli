# Marketplace CLI

`mkpcli` enables a command-line interface to the [VMware Marketplace](http://marketplace.cloud.vmware.com/) for consumers and publishes.

To install, grab the latest prebuilt binary from the [Releases](https://github.com/vmware-labs/marketplace-cli/releases) page, or [build from source](#building).

Features:
* Get details about a product
* Manage products in your org
  * Add versions
  * Attach container images
  * Attach Helm charts
  * Attach OVA files
* Download assets from a product

## Authentication

`mkpcli` requires an API Token from [VMware Cloud Services](https://console.cloud.vmware.com/csp/gateway/portal/#/user/tokens). See [this doc](./docs/Authentication.md) for more information.

## Example
<a href="https://asciinema.org/a/68HbJWxv13rmrOwukYhO72ndD" target="_blank">
  <img src="https://asciinema.org/a/68HbJWxv13rmrOwukYhO72ndD.svg" alt="Demo of mkpcli" />
</a>

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

Please see our [Code of Conduct](CODE-OF-CONDUCT.md) and [Contributors guide](CONTRIBUTING.md).

If you would like to build and contribute to this project, please fork and make pull requests. 
