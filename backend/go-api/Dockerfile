############################
# STEP 1 build executable binary
############################
FROM golang@sha256:dd8888bb7f1b0b05e1e625aa29483f50f38a9b64073a4db00b04076cec52b71c as builder
# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
# Create appuser
ENV USER=appuser
ENV UID=10001
# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .

# Fetch dependencies.
# Using go mod.
RUN go mod download
RUN go mod verify

# Build the binary
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/hello

############################
# STEP 2 build a small image
############################
FROM scratch
# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# Copy our static executable
COPY --from=builder /go/bin/hello /go/bin/hello
# Use an unprivileged user.
USER appuser:appuser
# Run the hello binary.
ENTRYPOINT ["/go/bin/hello"]
