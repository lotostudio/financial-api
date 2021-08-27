# Main dockerfile for building images on unix systems (build can be done manually)
#
# Generate swagger:
# swag init -g internal/app/app.go
#
# Build binary:
# go build -ldflags="-w -s" -v -o app ./cmd/main.go
#
# Build image:
# docker build -f build/Dockerfile -t github.com/lotostudio/financial-api .
#
# Run container by passing env variables by flag -e (list of env variables in README.md):
# docker run -p 8080:8080 github.com/lotostudio/financial-api:latest
#
# Transfer artifacts to clean alpine container
FROM alpine:latest

WORKDIR /root/

COPY . .

EXPOSE 8080

CMD ["./app"]
