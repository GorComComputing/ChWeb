package main

import (
	"io"

	"golang.org/x/net/websocket"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

// WebSockWrapper wraps the raw websocket and converts Write() calls
// to proper websocket.Send() working in binary or text mode. If text
// mode is selected, it converts the data passed to Write() into UTF8 bytes
//
// We need this to make sure that the entire buffer in io.Writer.Write(buffer)
// is delivered as a single chunk to the web browser, instead of being split
// into multiple frames. This wrapper basically substitues every Write() with
// Send() and every Read() with Receive()
type WebSockWrapper struct {
	io.ReadWriteCloser

	ws   *websocket.Conn
	mode WebSocketMode

	encoder *encoding.Encoder
	decoder *encoding.Decoder
}

// WebSocketMode allows to create WebSocket wrappers working in text
// or binary mode
type WebSocketMode int

const (
	WebSocketBinaryMode = iota
	WebSocketTextMode
)

func NewWebSockWrapper(ws *websocket.Conn, m WebSocketMode) *WebSockWrapper {
	if ws == nil {
		return nil
	}
	return &WebSockWrapper{
		ws:      ws,
		mode:    m,
		encoder: unicode.UTF8.NewEncoder(),
		decoder: unicode.UTF8.NewDecoder(),
	}
}

// Write implements io.WriteCloser for WebSockWriter (that's the reason we're
// wrapping the websocket)
//
// It replaces raw Write() with "Message.Send()"
func (w *WebSockWrapper) Write(data []byte) (n int, err error) {
	n = len(data)
	if w.mode == WebSocketBinaryMode {
		// binary send:
		err = websocket.Message.Send(w.ws, data)
		// text send:
	} else {
		var utf8 string
		utf8, err = w.encoder.String(string(data))
		err = websocket.Message.Send(w.ws, utf8)
	}
	if err != nil {
		n = 0
	}
	return n, err
}

// Read does the opposite of write: it replaces websocket's raw "Read" with
//
// It replaces raw Read() with "Message.Receive()"
func (w *WebSockWrapper) Read(out []byte) (n int, err error) {
	var data []byte

	if w.mode == WebSocketBinaryMode {
		err = websocket.Message.Receive(w.ws, &data)
	} else {
		var utf8 string
		err = websocket.Message.Receive(w.ws, &utf8)
		switch err {
		case nil:
			data, err = w.decoder.Bytes([]byte(utf8))
		case io.EOF:
			return 0, io.EOF
		}
	}
	if err != nil {
		return 0, err
	}
	return copy(out, data), nil
}

func (w *WebSockWrapper) Close() error {
	return w.ws.Close()
}
