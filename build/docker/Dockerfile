# -------- Stage 1: build --------
	FROM golang:1.22-alpine AS builder

	WORKDIR /src
	COPY go.mod go.sum ./
	RUN go mod download
	
	COPY . .
	RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /mappy ./cmd/server
	
	# -------- Stage 2: slim runtime --------
	FROM alpine:3.19
	
	WORKDIR /app
	COPY --from=builder /mappy /app/mappy
	COPY config/ ./config/
	
	EXPOSE 8080
	ENV APP_ENV=prod
	
	ENTRYPOINT ["./mappy"]
	