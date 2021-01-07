# crdoc

crdoc is a CLI to generate markdown documentation from Kubernetes `CustomResourceDefinition` YAML files.

## Install

```bash
go get github.com/roee88/crdoc
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

- There are no custom type information because this information is not available in the YAMLs.
- Supported apiVersions are `apiextensions.k8s.io/v1` and `apiextensions.k8s.io/v1beta1`.
