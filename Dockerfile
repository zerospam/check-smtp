ARG APP_PACKAGE="github.com/zerospam/check-smtp"
ARG APP_PATH="/go/src/${APP_PACKAGE}"
ARG APP_NAME="smtpChecker"

FROM zerospam/go-dep-docker as builder

ENV CGO_ENABLED=0
ENV GOOS=linux

ARG APP_PACKAGE
ARG APP_PATH
ARG APP_NAME

COPY . $APP_PATH
WORKDIR $APP_PATH

RUN dep ensure
RUN go build -a -installsuffix cgo -o $APP_NAME

FROM alpine:latest

ARG APP_PATH
ARG APP_NAME
RUN sed -i -e 's/dl-cdn/dl-4/' /etc/apk/repositories \
    && apk --update --no-cache add ca-certificates

COPY --from=builder ${APP_PATH}/${APP_NAME} /${APP_NAME}

CMD ["/smtpChecker"]
