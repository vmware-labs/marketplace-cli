# Publishing Container Image Products
The Marketplace CLI can update products that host container images.
The container image must be referenced by a publicly accessible repository, tag, and tag type.
The tag type is either `FIXED` or `FLOATING`

## Example
To do this, you can use the `mkpcli attach image` command:

```bash
mkpcli attach image --product hyperspace-database --version 1.0.1 --image-repository astrowidgets/hyperspacedb --tag 1.0.1 --tag-type FIXED --instructions 'docker run astrowidgets/hyperspacedb:1.0.1'
```

If this version is a new version for the product, pass the `--create-version` flag:

```bash
mkpcli attach image --product hyperspace-database --version 1.0.1 --create-version --image-repository astrowidgets/hyperspacedb --tag 1.0.1 --tag-type FIXED --instructions 'docker run astrowidgets/hyperspacedb:1.0.1'
```
