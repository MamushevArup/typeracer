# Stage 1: Build the Go application
FROM golang:latest as builder

WORKDIR /race

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /race/cmd/main/app cmd/main/app.go

# Stage 2: Build the final image
FROM scratch

WORKDIR /app

# Copy the binary from the first stage
COPY --from=builder /race/cmd/main/app /app/app

# Copy .env file
COPY --from=builder /race/.env /app/.env

# Copy .yml files from ./config directory
COPY --from=builder /race/config/config.yml /app/config/config.yml

# Copy .sql files from ./schema directory
COPY --from=builder /race/schema/ /app/schema/

# Command to run the executable
CMD ["/app/app"]