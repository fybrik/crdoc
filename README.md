# crdoc

Generate markdown documentation from Kubernetes `CustomResourceDefinition` YAML files.

## Install

Download the appropriate version for your platform from [Releases](https://github.com/mesh-for-data/crdoc/releases/latest).
You may want to install the binary to somewhere in your system's PATH such as `/usr/local/bin`.

Alternatively, if you have go 1.16 or later then you can also use `go install`. 
This will put the latest released version of `crdoc` in `$(go env GOPATH)/bin`:

```bash
go install github.com/mesh-for-data/crdoc@latest
```

## Usage

```bash
Output markdown documentation from Kubernetes CustomResourceDefinition YAML files

Usage:
  crdoc [flags]

Examples:

  # Generate example/output.md from example/crds using the default markdown.tmpl tempalte: 
  crdoc --resources example/crds --output example/output.md

  # Override template (builtin or custom):
  crdoc --resources example/crds --output example/output.md --template frontmatter.tmpl
  crdoc --resources example/crds --output example/output.md --template templates_folder/file.tmpl

  # Use a Table of Contents to filter and order CRDs
  crdoc --resources example/crds --output example/output.md --toc example/toc.yaml


Flags:
  -h, --help               help for crdoc
  -o, --output string      Path to output markdown file (required)
  -r, --resources string   Path to directory with CustomResourceDefinition YAML files (required)
  -t, --template string    Path to file in a templates directory (default "markdown.tmpl")
  -c, --toc string         Path to table of contents YAML file
```

## Limitations

- There are no custom type information because this information is not available in the YAMLs. This tool was specifically designed with that in mind and provides a different reader experience compared to other similar tools.
- `apiextensions.k8s.io/v1beta1` is supported but we assume a structural schema as required by v1.

## Similar tools

- [gen-resourcesdocs](https://github.com/kubernetes-sigs/reference-docs/tree/master/gen-resourcesdocs)
- [gen-crd-api-reference-docs](https://github.com/ahmetb/gen-crd-api-reference-docs)
- [crd-ref-docs](https://github.com/elastic/crd-ref-docs)

