# crdoc

[![Go Report Card](https://goreportcard.com/badge/github.com/fybrik/crdoc)](https://goreportcard.com/report/github.com/fybrik/crdoc)
[![Go Reference](https://pkg.go.dev/badge/github.com/fybrik/crdoc.svg)](https://pkg.go.dev/github.com/fybrik/crdoc)
[![golangci-lint](https://github.com/fybrik/crdoc/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/fybrik/crdoc/actions/workflows/golangci-lint.yml)
[![CodeQL](https://github.com/fybrik/crdoc/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/fybrik/crdoc/actions/workflows/codeql-analysis.yml)
[![gosec](https://github.com/fybrik/crdoc/actions/workflows/golang-security.yml/badge.svg)](https://github.com/fybrik/crdoc/actions/workflows/golang-security.yml)

Generate markdown documentation from Kubernetes `CustomResourceDefinition` YAML files.

## Install

Install `crdoc` using one of the following methods:


<details>
  <summary>Binary</summary>    

Download the appropriate version for your platform from [Releases](https://github.com/fybrik/crdoc/releases/latest). You may want to install the binary to somewhere in your system's PATH such as `/usr/local/bin`.

For convenience, an [installation script](https://raw.githubusercontent.com/fybrik/crdoc/main/hack/get-crdoc.sh) is also available:

</details>

<details>
  <summary>Go install</summary>    

If you have goLang 1.19 or later then you can use `go install`. 
This will put the latest released version of `crdoc` in `$(go env GOPATH)/bin`:

```bash
go install fybrik.io/crdoc@latest
```

> :bulb: Prefer pinning to a specific version rather than using @latest when installing in a CI workflow

</details>

<details>
  <summary>Docker images</summary>    

[Docker images](https://github.com/fybrik/crdoc/pkgs/container/crdoc) are available and allow you to run `crdoc` directly via docker. For example:

```bash
docker run -u $(id -u):$(id -g) --rm -v ${PWD}:/workdir ghcr.io/fybrik/crdoc:latest --resources /workdir/example/crds --output /workdir/example/output.md
```

> :bulb: Prefer pinning to a specific version rather than using :latest when running in a CI workflow

</details>

## Usage

```bash
Output markdown documentation from Kubernetes CustomResourceDefinition YAML files

Usage:
  crdoc [flags]

Examples:

  # Generate example/output.md from example/crds using the default markdown.tmpl template: 
  crdoc --resources example/crds --output example/output.md

  # Override template (builtin or custom):
  crdoc --resources example/crds --output example/output.md --template frontmatter.tmpl
  crdoc --resources example/crds --output example/output.md --template templates_folder/file.tmpl

  # Use a Table of Contents to filter and order CRDs
  crdoc --resources example/crds --output example/output.md --toc example/toc.yaml


Flags:
  -h, --help               help for crdoc
  -o, --output string      Path to output markdown file (required)
  -r, --resources string   Path to YAML file or directory containing CustomResourceDefinitions (required)
  -t, --template string    Path to file in a templates directory (default "markdown.tmpl")
  -c, --toc string         Path to table of contents YAML file
  -v, --version            version for crdoc
```

## Limitations

- There are no custom type information because this information is not available in the YAMLs. This tool was specifically designed with that in mind and provides a different reader experience compared to other similar tools.
- `apiextensions.k8s.io/v1beta1` is supported but we assume a structural schema as required by v1.

## Similar tools

- [gen-resourcesdocs](https://github.com/kubernetes-sigs/reference-docs/tree/master/gen-resourcesdocs)
- [gen-crd-api-reference-docs](https://github.com/ahmetb/gen-crd-api-reference-docs)
- [crd-ref-docs](https://github.com/elastic/crd-ref-docs)

