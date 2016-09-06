SET GOPATH=%CD%
rem go get github.com/gorilla/websocket
rem go get golang.org/x/net/websocket

go install code.naumchevski.com/iot/iotserver
go install code.naumchevski.com/iot/tester