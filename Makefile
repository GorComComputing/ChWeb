.PHONY: main
main: *.go deps
	GOOS=linux GOARCH=arm go build -o ChWeb .


.PHONY:deps
deps:
	go get github.com/gorilla/sessions
	go get github.com/kabukky/httpscerts
	go get github.com/lib/pq

	go get github.com/kr/pty
	go get golang.org/x/net/websocket
	go get golang.org/x/text/encoding
	go get golang.org/x/text/encoding/unicode

