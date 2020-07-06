![Latest Release Version](https://img.shields.io/github/v/release/cinarmert/klogs)
[![Go Report Card](https://goreportcard.com/badge/github.com/cinarmert/klogs)](https://goreportcard.com/report/github.com/cinarmert/klogs)
![Go CI](https://github.com/cinarmert/klogs/workflows/Go%20CI/badge.svg)

# `klogs`

**`klogs`** helps you view logs from multiple kubernetes pods and containers inside the pods. You can provide a regex to filter the pods/containers that the logs will be fetched from.

![klogs demo gif](img/klogs-demo.gif)

**`klogs`** is a minimalistic cli tool. See the usage below!

```
Klogs provides a user interface for the logs from multiple Kubernetes pods/containers.

Usage:
  klogs POD-QUERY [flags]

Flags:
  -c, --container string    Regex to specify the containers to fetch logs from (default ".*")
  -x, --context string      Kubernetes context to use, default is $(kubectl config current-context)
  -h, --help                help for klogs
  -k, --kubeconfig string   Kubeconfig file path
  -n, --namespace string    Kubernetes namespace to use, uses current namespace if not specified 
  -v, --verbose             print debug logs
```

# Installation

**`klogs`** is available on macOS, Linux and Windows. You can find the binaries in [**Releases &rarr;**](https://github.com/cinarmert/klogs/releases).

## Brew

**`klogs`** is already available in Homebrew for an easier installation.

```
brew tap cinarmert/klogs
brew install klogs
```

# Similar Projects

- [doclogs](https://github.com/cinarmert/doclogs)
- [stern](https://github.com/wercker/stern)
