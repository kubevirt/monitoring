FROM golang:1.20 as build

WORKDIR /workspace

COPY . .

RUN go mod tidy && \
    CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w" -o /bin/ .

FROM gcr.io/distroless/base:latest

COPY --from=build /bin/prom-metrics-linter .

ENTRYPOINT ["/prom-metrics-linter"]