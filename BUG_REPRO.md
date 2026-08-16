# Reproduction

Create the same problematic filename in two different directories so both files receive the same suggested replacement name. The dry-run plan reports only one of the two files.

Expected: paths in different directories are independent and both rename proposals are shown.

Command:

```bash
go run ./cmd/filename-normalizer --dir ./fixture --dry-run
```
