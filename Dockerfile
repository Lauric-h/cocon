FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY ui ./ui

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/cocon ./cmd/web


FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -H -u 10001 cocon \
    && mkdir /data \
    && chown cocon:cocon /data

WORKDIR /app

COPY --from=build /out/cocon ./cocon
COPY --from=build /src/ui ./ui

USER cocon

ENV DB_PATH=/data/cocon.db \
    PORT=8081 \
    TZ=Europe/Paris \
    GOOGLE_CALENDAR_ICS_URL=""

EXPOSE 8081

ENTRYPOINT ["/app/cocon"]