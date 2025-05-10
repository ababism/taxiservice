FROM golang:1.23.0 as build
WORKDIR /app

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o musicsnap-svc ./services/musicsnap/cmd

FROM alpine:latest as production

COPY --from=build /app/musicsnap-svc ./

COPY --from=build /app/.docker.env ./.env
COPY --from=build /app/services/musicsnap/migrations ./migrations
COPY --from=build /app/services/musicsnap/config/config.docker.yml ./services/musicsnap/config/config.local.yml

CMD ["./musicsnap-svc"]

EXPOSE 8080



#RUN apt-get update && \
#    apt-get --yes --no-install-recommends install make="4.3-4.1" && \
#    apt-get clean && rm -rf /var/lib/apt/lists/*