FROM golang:1.21-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
COPY busyworker ./busyworker/
RUN go mod download -x && go mod verify
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o busyworker/main busyworker/main.go

FROM alpine:latest
# RUN apk --no-cache add ca-certificates
COPY --from=build /src/busyworker/main /bin/busyworker