# Stage 1: Build the Go application
FROM golang:latest as builder

WORKDIR /race

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /race/cmd/main/app cmd/main/app.go

FROM scratch

WORKDIR /app

COPY --from=builder /race/cmd/main/app /app/app

COPY --from=builder /race/.env /app/.env

COPY --from=builder /race/config/config.yml /app/config/config.yml

COPY --from=builder /race/schema/ /app/schema/

CMD ["/app/app"]