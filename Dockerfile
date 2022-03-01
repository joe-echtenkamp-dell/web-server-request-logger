# syntax=docker/dockerfile:1

FROM golang:alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o ./proxy.out

##
## Deploy
##
FROM alpine
RUN apk update 

WORKDIR /

COPY --from=build /app/proxy.out /proxy.out

EXPOSE 8081

RUN addgroup -S nonroot && adduser -S nonroot -G nonroot 
USER nonroot

ENTRYPOINT ["/proxy.out"]
