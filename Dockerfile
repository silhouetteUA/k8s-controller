FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
ARG VERSION
ARG COMMIT
ARG DATE
RUN CGO_ENABLED=0 \
    GOOS=${GOOS} \
    GOARCH=${GOARCH} \
    go build -v -o k8s-controller \
    -ldflags "-X=github.com/silhouetteUA/k8s-controller/cmd.Version=${VERSION} \
              -X=github.com/silhouetteUA/k8s-controller/cmd.Commit=${COMMIT} \
              -X=github.com/silhouetteUA/k8s-controller/cmd.BuildDate=${DATE}" \
    main.go

# Final stage
FROM gcr.io/distroless/static-debian12
WORKDIR /
COPY --from=builder /app/k8s-controller .
EXPOSE 8080
ENTRYPOINT ["/k8s-controller"]
CMD ["server"]