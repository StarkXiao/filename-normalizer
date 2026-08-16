# Reproduction

Create a filename containing multibyte Unicode characters and set a maximum length measured in characters. A name within the configured character limit is reported as too long.

Expected: filename length is counted in Unicode characters, not encoded bytes.

Command:

```bash
go run ./cmd/filename-normalizer --dir ./fixture --max-length 12
```
