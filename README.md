# crdoc

crdoc is a CLI to generate markdown documentation from Kubernetes `CustomResourceDefinition` YAML files.

## Install

```bash
go get github.com/roee88/crdoc
```

This will put crdoc in `$(go env GOPATH)/bin`. You may need to add that directory to your `$PATH` if you encounter a "command not found" error.

## Usage

```
crdoc [flags]

Flags:
  -c, --config-dir string      Configuration directory (required)
  -h, --help                   help for crdoc
  -o, --output-dir string      Output directory for markdown file (required)
  -r, --resources-dir string   Directory containing CustomResourceDefinition YAML files (required)
  -t, --toc-file string        Path to table of contents YAML file
  
Required flags: config-dir, output-dir, resources-dir
```

For example:

```bash
mkdir -p example/output/
crdoc --config-dir config --resources-dir example/crds --output-dir example/output
```

## Limitations

- Supported apiVersions are `apiextensions.k8s.io/v1` and `apiextensions.k8s.io/v1beta1`.
- There are no custom type information because this information is not available in the YAMLs.
- Currently a single markdown file is generated name `out.md`.
