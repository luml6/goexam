FROM golang
MAINTAINER  albert.lu@zeusis.com
ENV PROJECTNAME notepad-api
ADD gopkg.in.tar.gz  $GOPATH/src/
ADD golang.org.tar.gz $GOPATH/src/
ADD github.com.tar.gz $GOPATH/src/

RUN cd $GOPATH/src && git clone  http://10.0.12.104/albert.lu/notepad-api.git
RUN cd $GOPATH/src/$PROJECTNAME && go get -d -v
RUN cd $GOPATH/src/$PROJECTNAME && go build
CMD cd $GOPATH/src/$PROJECTNAME &&  ./$PROJECTNAME

EXPOSE 8088 22

