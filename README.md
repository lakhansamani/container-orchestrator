# cotainer-orchestrator

gRPC Service to spin up and down a container.

## Getting started

### Required services

- Redis to maintain state of container

### Required envs

```sh
export REDIS_URL=redis://localhost:6379
```

### Starting server locally

- `make run`

### Creating Builds

- `make`
