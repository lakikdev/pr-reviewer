# Builder Stage
FROM golang:1.23.8 AS builder
WORKDIR /app

ARG APP_VERSION
ARG CGO_ENABLED=0

# Download dependencies
COPY go.mod . 
COPY go.sum .
RUN go mod download

# Copy source code and build the binary
COPY . .
RUN go build -o /server ./cmd/server

FROM golang:1.21.0 AS swaggerui
RUN curl -sL https://github.com/swagger-api/swagger-ui/archive/v3.20.5.tar.gz -o /tmp/swagger.tar.gz
RUN mkdir /tmp/swagger
RUN tar xvvzf /tmp/swagger.tar.gz -C /tmp/swagger --strip 1
RUN sed -i -e "s+https://petstore.swagger.io/v2/swagger.json+/api/v1/swagger.json+g" /tmp/swagger/dist/index.html

# Final Distroless Stage
FROM gcr.io/distroless/static-debian12
WORKDIR /app

# Copy required runtime files
COPY --from=builder /server /server
COPY --from=swaggerui /tmp/swagger/dist /swaggerui
COPY internal/database/migrations /migrations
COPY swagger.json /swagger.json

# Timezone support
COPY --from=builder /usr/share/zoneinfo/ /usr/share/zoneinfo/
ENV TZ=UTC

ENTRYPOINT ["/server"]