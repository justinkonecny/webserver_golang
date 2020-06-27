FROM golang:1.14.3

EXPOSE 8081

WORKDIR /webserver

RUN apt-get update
RUN apt-get -y upgrade

RUN go get github.com/gorilla/sessions
RUN go get github.com/gorilla/mux
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/jinzhu/gorm/dialects/mysql
RUN go get github.com/jinzhu/gorm
RUN go get github.com/aws/aws-sdk-go-v2/aws
RUN go get github.com/aws/aws-sdk-go-v2/aws/external
RUN go get github.com/aws/aws-sdk-go-v2/service/sns

ADD main /webserver/main
ADD server /webserver/server
ADD fullchain.pem /webserver
ADD privkey.pem /webserver

RUN go build -o /webserver/run_server /webserver/main/main.go

CMD /webserver/run_server