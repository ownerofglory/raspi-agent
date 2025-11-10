# UI build phase
FROM node:22-alpine AS ui-build
WORKDIR /ui

COPY ui/package*.json ./
RUN npm ci --ignore-scripts

COPY ui/ .
RUN npm run build

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

COPY --from=ui-build /ui/dist /app/ui/dist

# Run phase
FROM alpine:latest
RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot \

USER nonroot

WORKDIR /app

RUN apk --no-cache add ca-certificates

# Copy binary and static assets
COPY --from=build /app/bin/raspi-agent-backend .
COPY --from=build /app/ui/dist ./ui/dist

EXPOSE 8080
CMD ["./raspi-agent-backend"]