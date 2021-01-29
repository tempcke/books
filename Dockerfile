FROM golang:alpine as build
RUN  apk update \
  && apk add curl \
            git \
            bash \
            make \
            ca-certificates \
  && rm -rf /var/cache/apk/*

# install migrate which will be used by entrypoint.sh to perform DB migration
ARG MIGRATE_VERSION=4.7.1
ADD https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz /tmp
RUN  tar -xzf /tmp/migrate.linux-amd64.tar.gz -C /usr/local/bin \
  && mv /usr/local/bin/migrate.linux-amd64 /usr/local/bin/migrate

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN	make build


FROM alpine
RUN  apk --no-cache add ca-certificates bash \
  && update-ca-certificates \
  && mkdir -p /var/log/app
WORKDIR /app
COPY --from=build /usr/local/bin/migrate /usr/local/bin
COPY --from=build /app/db/migrations /app/db/migrations
COPY --from=build /app/bin /app/bin
EXPOSE 8080 8080
CMD ["/app/bin/bookserver"]