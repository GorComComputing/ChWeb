package main

import (
	"net/http"
    	"os/exec"
    	"log"
	"io"
	"github.com/kr/pty"
	"golang.org/x/net/websocket"
)


// GET /shell handler
// Launches /bin/bash and starts serving it via the terminal
func onShell(w http.ResponseWriter, r *http.Request) {
	wsHandler := func(ws *websocket.Conn) {
		// wrap the websocket into UTF-8 wrappers:
		wrapper := NewWebSockWrapper(ws, WebSocketTextMode)
		stdout := wrapper
		stderr := wrapper

		// this one is optional (solves some weird issues with vim running under shell)
		stdin := &InputWrapper{ws}

		// starts new command in a newly allocated terminal:
		// TODO: replace /bin/bash with:
		//		 kubectl exec -ti <pod> --container <container name> -- /bin/bash
		cmd := exec.Command("/bin/sh") //   /root/ChShell
		tty, err := pty.Start(cmd)
		if err != nil {
			panic(err)
		}
		defer tty.Close()

		// pipe to/fro websocket to the TTY:
		go func() {
			io.Copy(stdout, tty)
		}()
		go func() {
			io.Copy(stderr, tty)
		}()
		go func() {
			io.Copy(tty, stdin)
		}()

		// wait for the command to exit, then close the websocket
		cmd.Wait()
	}
	defer log.Printf("Websocket session closed for %v", r.RemoteAddr)

	// start the websocket session:
	websocket.Handler(wsHandler).ServeHTTP(w, r)
}
