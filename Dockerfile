# Stage 1: Build the app
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/backdev_go .


# # Stage 2: Create the final runtime image
FROM alpine:3.19 AS final
WORKDIR /app
COPY ./static /app/static
COPY ./config.toml /app/config.toml
COPY --from=builder ./app/backdev_go /app/

# Expose any required port (optional)
EXPOSE 8080

# Run the Go executable
ENTRYPOINT ["./backdev_go"]
CMD ["config.toml"]