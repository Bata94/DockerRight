### BASE
FROM golang:1.22-bookworm as base

WORKDIR /opt/DockerRight

COPY . .

RUN go mod tidy

### PRODUCTION-BUILDER
FROM base as prod-builder

WORKDIR /opt/DockerRight

COPY --from=base /opt/DockerRight .

RUN go build -o ./DockerRight cmd/main.go

### PRODUCTION
FROM golang:1.22-bookworm as prod

WORKDIR /opt/DockerRight

COPY --from=prod-builder /opt/DockerRight/DockerRight .

ENTRYPOINT ["/opt/DockerRight/DockerRight"]

### DEVELOPMENT 
FROM base as dev

WORKDIR /opt/DockerRight

COPY --from=base /opt/DockerRight .

CMD ["go", "run", "cmd/main.go"]
