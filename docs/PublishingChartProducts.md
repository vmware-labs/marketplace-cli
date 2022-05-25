# Publishing Chart Products
The Marketplace CLI can update products that host chart images.
The CLI can either upload a local chart (in a directory or tgz format), or attach a chart based on a public URL, then it
will attach the reference to the product.

## Example
To do this, you can use the `mkpcli attach chart` command:

### Upload a local chart

```bash
mkpcli attach chart --product hyperspace-database-chart --product-version 1.0.1 --chart charts/hyperspace-db-1.0.1.tgz --instructions 'helm install it'
```

### Attaching a remote chart

```bash
mkpcli attach chart --product hyperspace-database-chart --product-version 1.0.1 --chart https://astro-widgets.example.com/charts/hyperspace-db-1.0.1.tgz --instructions 'helm install it'
```

If this version is a new version for the product, pass the `--create-version` flag:

```bash
mkpcli attach chart --product hyperspace-database-chart --product-version 1.0.1 --create-version --chart charts/hyperspace-db-1.0.1.tgz --instructions 'helm install it'
```
