# --- build stage ---
FROM golang:1.26 AS build

WORKDIR /src

# Cache deps first.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/lnkd ./cmd/lnkd

# --- runtime stage ---
FROM gcr.io/distroless/static-debian12

COPY --from=build /out/lnkd /usr/local/bin/lnkd

EXPOSE 8080
ENTRYPOINT ["lnkd"]
