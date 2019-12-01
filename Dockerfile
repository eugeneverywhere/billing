FROM golang:1.12 as build_base

COPY .netrc /root/.netrc

# build envs
ENV GOOS linux
ENV GOARCH amd64

# service configs
ENV SERVICE_NAME billing
ENV GIT_SSL_NO_VERIFY 1
ENV PKG github.com/eugeneverywhere/billing

RUN mkdir -p /go/src/${PKG}
WORKDIR /go/src/${PKG}/${SERVICE_NAME}

COPY go.mod go.sum Makefile ./



FROM build_base as service_builder

COPY . .


FROM alpine:3.7 as base_runner

ENV PKG github.com/eugeneverywhere/billing

ENV PAYMENT_PROCESSOR_NAME   				billing


WORKDIR /root/
RUN apk add curl


FROM base_runner as billing
COPY --from=service_builder /go/src/${PKG}/cmd/${PAYMENT_PROCESSOR_NAME}/${PAYMENT_PROCESSOR_NAME} .
ENTRYPOINT ["./billing", "-config=local.yml"]

FROM service_builder as service_tester
RUN make test