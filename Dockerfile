FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /site .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /site /site
EXPOSE 8443/tcp 8443/udp
ENTRYPOINT ["/site"]
