FROM jamesgroat/golang-nodejs

RUN npm install -g bower
ADD package.json /keepthestreakalive/package.json
WORKDIR /keepthestreakalive
RUN npm install

ADD public/bower.json /keepthestreakalive/public/bower.json
WORKDIR /keepthestreakalive/public
RUN bower install

ADD main.go /keepthestreakalive/main.go
WORKDIR /keepthestreakalive
RUN go build main.go

ADD public/ /keepthestreakalive/public/
WORKDIR /keepthestreakalive/public
RUN npm run build

WORKDIR /keepthestreakalive
CMD ["/keepthestreakalive/main"]


