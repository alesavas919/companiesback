# Image build
FROM golang:1.24.4-alpine AS builder

# Stablish work dir path into container
WORKDIR /app

COPY go.mod go.sum ./

# install go depdencies
RUN go mod tidy

# Copy all files from container
COPY . .

RUN go get -u all

RUN go get github.com/aws/aws-sdk-go-v2

# Binary go build
RUN go build -o alejo101_companies .

# Image excecute
FROM alpine:latest

# Install necesary certicates
RUN apk --no-cache add ca-certificates

# Copy compile binary until the build tie
COPY --from=builder /app/alejo101_companies /usr/local/bin/alejo101_companies

# Port application expose
EXPOSE 8081

# ARGS
ARG SSL_TOKEN_PROJECT
ARG SSL_SECRET_SL
ARG SSL_ADD_PROJECT
ARG SSL_REGION
ARG TEST_ENV_A

ENV SSL_TOKEN_PROJECT=${SSL_TOKEN_PROJECT}
ENV SSL_SECRET_SL=${SSL_SECRET_SL}
ENV SSL_ADD_PROJECT=${SSL_ADD_PROJECT}
ENV SSL_REGION=${SSL_REGION}
ENV TEST_ENV_A=${TEST_ENV_A}



# Run app command
CMD ["/usr/local/bin/alejo101_companies"]