# Repository

## Build

```bash
helm package .
helm repo index .
```

## Add

```bash
helm repo add check-smtp https://raw.githubusercontent.com/zerospam/check-smtp/helm-chart/
```

## Install

```bash
helm install check-smtp/check-smtp
```