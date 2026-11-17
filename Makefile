.PHONY: test run

# Same checks as the PR workflow's test job.
test:
	go vet ./...
	go test -race ./...

# Runs the service with the settings in .env (copy .env.example first).
run:
	go run ./cmd/agneta
