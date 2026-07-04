# ---- build ----
FROM golang:1.25 AS build
WORKDIR /src

# Baixa dependências primeiro (cache de camada).
COPY go.mod go.sum ./
RUN go mod download

# Compila o binário estático (sem CGO) — roda no distroless.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/patient-data-service ./cmd/server

# ---- runtime ----
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/patient-data-service /patient-data-service
EXPOSE 50051 9090
USER nonroot:nonroot
ENTRYPOINT ["/patient-data-service"]
