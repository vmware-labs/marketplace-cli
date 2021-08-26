# Marketplace CLI

The Marketplace CLI is a tool that can be used to interact with the [VMware Marketplace](http://marketplace.cloud.vmware.com/).
The primary focus for the CLI is to assist publishers with automation.

## Examples

### CI/CI examples

Adding a new version using [Concourse](https://concourse-ci.org/):
```yaml
resources:
- name: mkpcli
  type: docker-image
  source:
    repository: harbor-repo.vmware.com/tanzu_isv_engineering/mkpcli
- name: version
  type: semver
  source: ...

jobs:
- name: Add version
  plan:
  - get: mkpcli
  - get: version
    params: { bump: patch }
  - task: add-version-to-marketplace
    image: mkpcli
    config:
      inputs:
        - name: version
      platform: linux
      params:
        CSP_API_TOKEN: ((marketplace_api_token))
        SLUG: test-container-product2
      run:
        path: bash
        args:
        - -exc
        - |
          mkpcli product-version create \
            --product my-marketplace-product1 \
            --product-version $(cat version/version)
```

## Building

Building from source is simple with our Makefile.

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

Development [Roadmap](https://miro.com/app/board/o9J_l_2uPFI=/)

Please see our [Code of Conduct](CODE-OF-CONDUCT.md) and [Contributors guide](CONTRIBUTING.md).

A few prerequisites will be helpful for setting up your development environment:

### Set up vault (VMware internal)

Get, configure, and log in to vault:

```bash
$ brew install vault
...
$ export VAULT_ADDR=https://runway-vault.svc.eng.vmware.com
$ vault login -method=ldap username=<username>
Password (will be hidden):
Success! You are now authenticated.
...
```

### Enable direnv

Direnv allows for settings to be loaded when entering the directory. Here, it simplifies setting up the development environment so you don't forget.

```bash
brew install direnv
direnv allow
```
### Set GOPATH
Make sure GOPATH is set. e.g. export `PATH=$PATH:$(go env GOPATH)/bin`.
