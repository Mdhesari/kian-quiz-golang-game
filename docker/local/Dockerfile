# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang:1.12-alpine base image
FROM golang:latest

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
# RUN apk update && apk upgrade && \
#     apk add --no-cache bash git openssh

RUN --mount=type=cache,target=/var/cache/apt \
    apt-get update && apt-get install -y build-essential

WORKDIR /usr/src/app

COPY go.* .

RUN go mod tidy
RUN go mod verify
RUN go mod download

COPY . .