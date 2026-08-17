# Filename Normalizer

`filename-normalizer` scans a directory for filename issues and can emit a CSV report or apply non-conflicting rename suggestions.

## Build

```bash
go build ./...
```

## Test

```bash
go test ./...
```

## Run

Scan a directory without changing files:

```bash
go run ./cmd/filename-normalizer -dir /path/to/files -dry-run
```

Write findings to CSV:

```bash
go run ./cmd/filename-normalizer -dir /path/to/files -csv findings.csv
```

Apply non-conflicting rename suggestions:

```bash
go run ./cmd/filename-normalizer -dir /path/to/files -rename
```

The command exits with status 1 when it finds filename issues.

## Docker

Build the evaluation image (defaults to `linux/amd64`):

```bash
./build_benzhi_docker.sh filename-normalizer
```

Build for Apple Silicon or x86 explicitly:

```bash
./build_benzhi_docker.sh filename-normalizer linux/arm64
./build_benzhi_docker.sh filename-normalizer linux/amd64
```

Start an interactive container:

```bash
docker run -it filename-normalizer:latest
```
