# ---- builder ----

FROM golang:1.23.3 AS builder

WORKDIR /app

# dependency caching
COPY go.mod go.sum ./

RUN go mod download

# source
COPY . .

# prod build: static, trimmed
RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-s -w" -trimpath -o students_api ./cmd/api


# ---- debug runtime(has shell) ----
  
FROM debian:bookworm-slim AS debug-run

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl bash procps && rm -rf /var/lib/apt/lists*
   

COPY --from=builder /app/students_api .

ENV PORT=8080

EXPOSE 8080

CMD ["./students_api"]

# ---- prod runtime (distroless) ----

FROM gcr.io/distroless/base-debian12 AS prod-run

WORKDIR /app

COPY --from=builder /app/students_api .

ENV PORT=8080

EXPOSE 8080

CMD ["./students_api"]


