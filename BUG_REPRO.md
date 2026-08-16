# Reproduction

Set a maximum filename length and create one filename whose length is exactly that value. The scanner reports the boundary-value filename as too long.

Expected: a name at the configured maximum is accepted; only a name exceeding the maximum is reported.

Command:

```bash
go run ./cmd/filename-normalizer --dir ./fixture --max-length 12
```
