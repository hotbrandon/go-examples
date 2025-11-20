package main

import (
	"fmt"
	"net/http"
	"time"
)

// HTTP is stateless, but there is a distinction between the Application Layer (HTTP) and the Transport Layer (TCP) that creates this "magic."

// Here is exactly how the server knows the user is gone.

// The "Stateless" Misunderstanding
// When we say "HTTP is stateless," we mean the server doesn't remember who you are between request A and request B (unless we use cookies/sessions).

// However, during a single request (Request A), there is an active, live TCP Connection (a socket) kept open between the browser and the server.

// The Chain of Events: How the "Stop" Signal Travels
// Imagine you are downloading a large file or waiting for a slow database query.

// The Browser: When you close the tab, the browser knows it no longer needs the response. It immediately sends a TCP FIN (Finish) or RST (Reset) packet to the server. It essentially "hangs up the phone."

// The OS/Kernel: Your server's operating system (Linux/Windows) receives this packet and realizes the socket is dead.

// Go net/http Server: The Go HTTP server is constantly monitoring that open socket. As soon as it detects the network connection is broken, it cancels the context.Context associated with that request.

// Your Code: Because you passed that ctx into your database query, the database driver sees the cancellation signal and stops executing the SQL.

// Let's Prove It (A Runnable Experiment)
// You can test this yourself with a tiny Go program. Run this code, open the page in your browser, and then immediately close the tab before the 5 seconds are up.

func main() {
	http.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("🟢 Request started...")

		// Use the request context
		ctx := r.Context()

		select {
		case <-time.After(5 * time.Second):
			// If 5 seconds pass, we finish successfully
			fmt.Println("✅ Request finished successfully")
			w.Write([]byte("Done!"))
		case <-ctx.Done():
			// If the browser closes the connection, this case triggers!
			err := ctx.Err()
			fmt.Println("❌ CLIENT DISCONNECTED! Stopping work:", err)
			// We return early, saving resources
			return
		}
	})

	fmt.Println("Server running on :8080. Go to http://localhost:8080/slow")
	http.ListenAndServe(":8080", nil)
}
