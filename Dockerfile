#-------------------------------------------
# STEP 1 : build executable binary
#-------------------------------------------
FROM golang:1.18.6-alpine3.16 as builder

ADD . /usr/src/app
WORKDIR /usr/src/app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager operator/main.go

#-------------------------------------------
# STEP 2 : build a image
#-------------------------------------------
FROM gcr.io/distroless/static:nonroot
#FROM busybox:latest

COPY --from=builder /usr/src/app/manager /app/bin/
COPY --from=builder /usr/src/app/manifests/charts /app/manifests/charts
COPY --from=builder /usr/src/app/manifests/profiles /app/manifests/profiles

#-------------------------------------------
# STEP 3 : execute a binary
#-------------------------------------------
USER 65532:6:5532
ENTRYPOINT [ "/app/bin/manager" ]
