FROM jamesgroat/golang-nodejs

RUN npm install -g bower
WORKDIR /keepthestreakalive
RUN npm install -g node-sass
ADD package.json /keepthestreakalive/package.json
WORKDIR /keepthestreakalive
RUN npm install

ADD public/bower.json /keepthestreakalive/public/bower.json
WORKDIR /keepthestreakalive/public
RUN bower install --allow-root

ADD main.go /keepthestreakalive/main.go
WORKDIR /keepthestreakalive
RUN go get github.com/gorilla/mux
RUN go install github.com/gorilla/mux
RUN go get github.com/yhat/scrape
RUN go install github.com/yhat/scrape
RUN go get golang.org/x/net/html
RUN go install golang.org/x/net/html
RUN go build main.go

ADD public/ /keepthestreakalive/public/
WORKDIR /keepthestreakalive/public
RUN node-sass ./scss/main.scss ./dist/main.css

EXPOSE 9898

WORKDIR /keepthestreakalive
CMD ["/keepthestreakalive/main"]


