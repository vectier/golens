# GoLens

**GoLens** is a language server that provides [CodeLens-style](https://code.visualstudio.com/blogs/2017/02/12/code-lens-roundup) implementation and references usage counts for [Go](https://go.dev/), similar to what you get in Rust, TypeScript, and other language.

This project exists as a workaround while waiting for native CodeLens support in `gopls` ([golang/go#56695](https://github.com/golang/go/issues/56695)).

> Built entirely with Go standrad libraries

## Features
- Show implementation counts for interface
- Show reference usage counts for interface

## Editor Integration

Currently available extensions:
- Zed - [zed-golens](https://github.com/riflowth/zed-golens)
