# Dependency Graph ፨

[godoclink]: https://godoc.org/github.com/manishmeganathan/depgraph
[![go docs](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godoclink]
![go version](https://img.shields.io/github/go-mod/go-version/manishmeganathan/depgraph?style=flat-square)
![latest tag](https://img.shields.io/github/v/tag/manishmeganathan/depgraph?color=brightgreen&label=latest%20tag&sort=semver&style=flat-square)
![license](https://img.shields.io/github/license/manishmeganathan/depgraph?color=g&style=flat-square)
![issue count](https://img.shields.io/github/issues/manishmeganathan/symbolizer?style=flat-square&color=yellow)

A Go Implementation for a simple dependency graph with circular resolution

### Overview
This package provides a simple implementation for a dependency graph with `DependencyGraph` into which edges and vertices can
be inserted, removed, checked for inclusivity, and iterated upon. It supports circular dependency resolution and deep fetching 
of dependencies. `DependencyGraph` can also be serialized into JSON, YAML, and POLO formats and implements the `engineio.DepDriver` interface for `go-moi`

### Installation
```
go get github.com/manishmeganathan/depgraph
```

### Notes:
This package is still a work in progress and can be heavily extended for a lot of different use cases.
If you are using this package and need some new functionality, please open an issue or a pull request.
