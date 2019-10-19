# Shell config file, included by some watcher (build on top of reflex)

# simple watch
dev_watch() {
  go test -coverprofile=cover.out ./... \
  && go tool cover -html=cover.out -o coverage.html \
  && reflex -R '^goflow' -r '\.go$' -- \
    sh -c 'go test -coverprofile=cover.out ./... && go tool cover -html=cover.out -o coverage.html'
}

# watch with Delve (Go debugger)
dev_watch_debug() {
  reflex --start-service -g '*.go' -- sh -c \
  "dlv debug $dev_entrypoint --headless=true --listen=0.0.0.0:2345 --api-version=2"
}