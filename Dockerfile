FROM golang:1.18-alpine3.16 as builder
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY ./cmd ./cmd
COPY ./internal ./internal
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64  \
    go build -ldflags '-w -s -extldflags "-static"' -o /app ./cmd/plgen

FROM scratch
COPY --from=builder /app /opt/backend
COPY ./templates /opt/templates
EXPOSE 8081/tcp
ENTRYPOINT ["/opt/backend", "-c", "/opt/config.yml"]
