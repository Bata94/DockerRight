ARG VERSION_BUILD=${VERSION_BUILD}

### BASE
FROM golang:1.22-bookworm as base

ARG VERSION_BUILD
ENV VERSION=$VERSION_BUILD

WORKDIR /opt/DockerRight

COPY . .

RUN go mod tidy

### PRODUCTION-BUILDER
FROM base as prod-builder

ARG VERSION_BUILD
ENV VERSION=$VERSION_BUILD

WORKDIR /opt/DockerRight

COPY --from=base /opt/DockerRight .

RUN go build -o ./DockerRight cmd/main.go

### PRODUCTION
FROM golang:1.22-bookworm as prod

ARG VERSION_BUILD
ENV VERSION=$VERSION_BUILD

WORKDIR /opt/DockerRight

COPY --from=prod-builder /opt/DockerRight/DockerRight .

ENTRYPOINT ["/opt/DockerRight/DockerRight"]

### DEVELOPMENT 
FROM base as dev

ENV VERSION="Development"

WORKDIR /opt/DockerRight

COPY --from=base /opt/DockerRight .

CMD ["go", "run", "cmd/main.go"]
