FROM golang:1.14.3

EXPOSE 8443

WORKDIR /webserver

RUN apt-get update
RUN apt-get upgrade -y

RUN go get github.com/gorilla/sessions
RUN go get github.com/gorilla/mux
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/jinzhu/gorm/dialects/mysql
RUN go get github.com/jinzhu/gorm
RUN go get github.com/google/uuid
RUN go get golang.org/x/crypto/bcrypt

ADD calendays /webserver/calendays
ADD common /webserver/common
ADD libertycars /webserver/libertycars
ADD main /webserver/main
ADD q /webserver/q

RUN go build -o /webserver/run_server /webserver/main/main.go

CMD /webserver/run_server
