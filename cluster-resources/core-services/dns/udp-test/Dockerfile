FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -o /bin/udp-echo

FROM alpine
COPY --from=build /bin/udp-echo /bin/udp-echo

EXPOSE 1234

ENTRYPOINT ["/bin/udp-echo", "server"]
