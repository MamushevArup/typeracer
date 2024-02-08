# Stage 1: Build the Go application
FROM golang:latest as builder

WORKDIR /race

COPY go.mod go.sum ./

RUN go mod download

COPY config/config.yml race/config.yml

COPY schema/*.sql race/schema/

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /race/cmd/main/app cmd/main/app.go

# Stage 2: Build the final image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the first stage
COPY --from=builder /race/cmd/main/app /app/app

# Copy .sql files from ./schema directory
COPY --from=builder /race/schema/ ./schema/

# Copy .yml files from ./config directory
COPY --from=builder /race/config/config.yml ./config/config.yml

# Command to run the executable
CMD ["/app/app"]
