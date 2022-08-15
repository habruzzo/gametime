FROM golang:1.18-alpine as build
COPY . /src
WORKDIR /src/gametime
RUN apk --no-cache add ca-certificates \
  && update-ca-certificates
RUN go mod download
RUN go build -o ./bin/gametime ./cmd/gametime

FROM alpine
COPY --from=build /src /src
WORKDIR /src/gametime
CMD ["./bin/gametime"]