# syntax=docker/dockerfile:1

FROM ghcr.io/merliot/device

WORKDIR /app
COPY . .

RUN go generate ./...
RUN go build -tags prime -o /garage ./cmd

EXPOSE 8000

ENV PORT_PRIME=8000
CMD ["/garage"]
