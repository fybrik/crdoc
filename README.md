# crdoc

crdoc is a CLI to generate markdown documentation from Kubernetes `CustomResourceDefinition` YAML files.

## Install

```bash
go get github.com/mesh-for-data/crdoc
```

This will put crdoc in `$(go env GOPATH)/bin`. You may need to add that directory to your `$PATH` if you encounter a "command not found" error.

## Usage

```
Output markdown documentation from Kubernetes CustomResourceDefinition YAML files

Usage:
  crdoc [flags]

Flags:
  -h, --help               help for crdoc
  -o, --output string      Path to output markdown file (required)
  -r, --resources string   Path to directory with CustomResourceDefinition YAML files (required)
  -t, --template string    Path to template file or directory (required)
  -c, --toc string         Path to table of contents YAML file
```

For example:

```bash
crdoc --template templates --resources example/crds --output out/output.md
```

## Limitations

- There are no custom type information because this information is not available in the YAMLs. This tool was specifically designed with that in mind and provides a different reader experience compared to other similar tools.
- `apiextensions.k8s.io/v1beta1` is supported but we assume a structural schema as required by v1.

## Similar tools

- [gen-resourcesdocs](https://github.com/kubernetes-sigs/reference-docs/tree/master/gen-resourcesdocs)
- [gen-crd-api-reference-docs](https://github.com/ahmetb/gen-crd-api-reference-docs)
- [crd-ref-docs](https://github.com/elastic/crd-ref-docs)

