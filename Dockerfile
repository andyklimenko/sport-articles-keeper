FROM golang:1.21-alpine as build

RUN apk add --no-cache make
ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.0
RUN make build

FROM alpine:3.18.3
COPY --from=build /app/articles-keeper /app/
ENTRYPOINT ["/app/articles-keeper"]
