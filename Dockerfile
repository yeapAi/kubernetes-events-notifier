FROM golang:1.17-alpine as build
ARG APPUSER=appuser

ENV USER=${APPUSER}
ENV UID=1001

RUN adduser -D -g "" -H -s "/sbin/nologin" -u "${UID}" "${USER}"
WORKDIR /app

COPY . ./

RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch
ARG APPUSER=appuser

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER ${APPUSER}:${APPUSER}
COPY --from=build --chown=${APPUSER}:${APPUSER} /app/main /go/bin/

WORKDIR /home/${APPUSER}
COPY config .kube/config

ENTRYPOINT ["/go/bin/main"]
