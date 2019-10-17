FROM golang

ARG app_env
ENV APP_ENV $app_env

COPY ./app /go/src/app
WORKDIR /go/src/app

RUN go get ./
RUN go build

CMD ["app"]
