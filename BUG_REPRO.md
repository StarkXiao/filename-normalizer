# Reproduction

Create a file inside a nested directory and run the checker with recursive scanning enabled. The nested file is omitted from the findings. With recursive scanning disabled, nested directories are incorrectly traversed.

Expected: recursive mode scans nested directories, and non-recursive mode scans only the selected directory.

Command:

```bash
go run ./cmd/filename-normalizer --dir ./fixture --recursive
```
