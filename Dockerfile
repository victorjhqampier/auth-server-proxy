# Etapa de construcción
FROM golang:1.20.4 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY src/ src/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./src/presentation

# Etapa de producción
FROM scratch
COPY --from=build /app/main /app/main
WORKDIR /app
CMD ["./main"]
