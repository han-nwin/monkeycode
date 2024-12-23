package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/han-nwin/monkeycode/tui"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
    "strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "localhost"
	port = "23234"
)

func main() {
	// Create the Wish SSH server
	server, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath("~/.ssh/id_ed25519"), // Path to your SSH host key
		wish.WithMiddleware(
            bubbletea.Middleware(teaHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		fmt.Printf("Error starting SSH server: %v\n", err)
		os.Exit(1)
	}

	// Graceful shutdown setup
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the SSH server
	fmt.Printf("SSH server started on %s:%s\n", host, port)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			fmt.Printf("Error running SSH server: %v\n", err)
			done <- nil
		}
	}()

	<-done
	fmt.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Error shutting down server: %v\n", err)
	}
}

// You can write your own custom bubbletea middleware that wraps tea.Program.
// Make sure you set the program input and output to ssh.Session.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
    // Load the typing prompt using the shared TUI function
    prompt := tui.LoadPrompt()

    // Create the Bubble Tea model
    model := tui.Model{
        PromptText:    prompt,
        UserText:      &strings.Builder{},
        CursorVisible: true,
        StartTime:     time.Now(),
        Width:         80, // Default terminal width; can be updated dynamically
        Height:        24, // Default terminal height
    }

    // Initialize the Bubble Tea program options
    options := []tea.ProgramOption{
        tea.WithAltScreen(),          // Use an alternate screen for better TUI experience
        tea.WithInput(s),             // Set input to the SSH session
        tea.WithOutput(s),            // Set output to the SSH session
    }

    // Return the initialized model and options
    return model, options
}


