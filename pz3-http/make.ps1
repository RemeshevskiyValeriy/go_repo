param(
    [Parameter(Mandatory=$true)]
    [ValidateSet("run", "build", "test")]
    [string]$Target
)

switch ($Target) {
    "run" {
        go run ./cmd/server
    }
    "build" {
        go build -o bin/server.exe ./cmd/server
    }
    "test" {
        go test ./internal/api
    }
}