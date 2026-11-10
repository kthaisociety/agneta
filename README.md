# agneta
Capture what advisors and the board know but never wrote down, so they can be offboarded and the board still has guidance.

---

## Development

The backend is Go 1.27 or later. Run the tests with:

```
go test ./...
```

Run the service with `go run ./cmd/agneta`. It listens on the port in `PORT`, or 8000 when unset, and answers `GET /health` with `ok`.

The interview schema lives in `internal/interview/schema.json`, with an example document next to it in `example.json`.
