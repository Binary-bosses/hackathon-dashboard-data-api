FROM golang:1.12-alpine3.9

WORKDIR /src

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Grabbing dependencies
COPY go.* ./

RUN go mod download

# Inlcude files
ADD . .

# Building app
RUN mkdir /app; go build -o goapp && cp goapp /app/

# Just the runtime, no build tools
FROM alpine:3.9
WORKDIR /app
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=0 /app /app
EXPOSE 8080
ENTRYPOINT ["./goapp"]