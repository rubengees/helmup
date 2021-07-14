# helmup

Check your helm chart dependencies for updates and apply them with an interactive selection. Helmup is a self contained
binary with zero dependencies written in go.

### Demo

[![asciicast](https://asciinema.org/a/iMMYijOOtJb2fjUBJK1Wms3Qf.svg)](https://asciinema.org/a/iMMYijOOtJb2fjUBJK1Wms3Qf)

### Usage

```shell
# Run in a directory with a Chart.yaml file.
helmup

# Pass a directory containing a Chart.yaml file.
helmup /path/to/helm/directory
```
