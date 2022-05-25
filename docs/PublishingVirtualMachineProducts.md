# Publishing Virtual Machine Products
The Marketplace CLI can update products that host virtual machine images (both ISO and OVA format).
The CLI will upload the image to the Marketplace, and then attach the reference to the product. 

## Example
To do this, you can use the `mkpcli attach vm` command:

```bash
mkpcli attach vm --product hyperspace-database-vm --product-version 1.0.1 --file vm/hyperspace-db-1.0.1-1526e30ba.iso
```

If this version is a new version for the product, pass the `--create-version` flag:

```bash
mkpcli attach vm --product hyperspace-database-vm --product-version 1.0.1 --create-version --file vm/hyperspace-db-1.0.1-1526e30ba.iso
```
