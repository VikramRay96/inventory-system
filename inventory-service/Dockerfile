FROM golang:1.19.1-alpine3.16
# Add Maintainer Info
VOLUME /var/log/
ENV LOG_PATH=/var/log/
ARG environment
ENV BOOT_CUR_ENV=${environment}
ENV AWS_KEY=${aws_key}
ENV AWS_SECRET=${aws_Secret}
# Set the Current Working Directory inside the container
WORKDIR /inventory-system

RUN apk add --no-cache git
ENV GOPRIVATE="bitbucket.org/kodnest"
RUN git config --global url."https://vikram_kod:ATBBqZdJKV5SgnRuvV2ppKCQeKZ706A29557@bitbucket.org".insteadOf "https://bitbucket.org"

# Copy go mod and sum files
COPY go.mod .
COPY go.sum .
RUN go clean -modcache
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY common common
COPY inventory-service inventory-service

# Build the Go app
RUN go build -o main ./inventory-service

# Expose port 8080 to the outside world
EXPOSE 8080

CMD  ["sh", "-c", "./main"]