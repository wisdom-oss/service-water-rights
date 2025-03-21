FROM docker.io/golang:alpine AS build-service
ENV GOMODCACHE=/root/.cache/go-build
WORKDIR /src
COPY --link go.* .
RUN --mount=type=cache,target=/root/.cache/go-build go mod download
COPY --link . .
RUN --mount=type=cache,target=/root/.cache/go-build go build -tags=release,nomsgpack,go_json -ldflags="-s -w" -o /service .

FROM --platform=$BUILDPLATFORM docker.io/alpine:latest AS compressor
COPY --from=build-service /service /service
RUN apk add --no-cache upx
RUN upx --best -o /compressed-service /service 


FROM scratch
ARG GH_REPO=unset
ARG GH_VERSION=unset
LABEL org.opencontainers.image.source=https://github.com/$GH_REPO
LABEL org.opencontainers.image.version=$GH_VERSION

COPY --from=build-service /etc/ssl/cert.pem /etc/ssl/cert.pem
COPY --from=compressor /compressed-service /service
ENTRYPOINT ["/service"]
EXPOSE 8000