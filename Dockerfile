# Build phase
FROM golang:1.24 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY Makefile ./
COPY . ./

ENV CGO_ENABLED=0
RUN make build-backend

# Run phase
FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/bin/raspi-agent-backend /usr/local/bin/raspi-agent-backend
EXPOSE 8080

CMD ["raspi-agent-backend"]