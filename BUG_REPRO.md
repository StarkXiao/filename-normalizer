# Reproduction

Place `report.txt` and `report.csv` in the same directory. The scanner reports a case conflict even though the files differ by extension and neither is a case-only duplicate.

Expected: only names that differ by letter case while keeping the same extension are grouped as a case conflict.

Command:

```bash
go run ./cmd/filename-normalizer --dir ./fixture
```
