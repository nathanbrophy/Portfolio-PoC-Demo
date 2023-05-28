FROM golang:1.20

COPY LICENSE /LICENSE

# Do not run as root
# We tie the new user to the root group to support arbitrary userID 
# file permission settings in some k8s distrobutions like OpenShift
RUN useradd -d /app -s /bin/bash -g root -G sudo -u 1001 serveruser
USER serveruser

WORKDIR /app

COPY . /app

# Download Go modules
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARH=amd64 go build -a -o /app/server


EXPOSE 8081

# Run
CMD ["/app/server"]