package server

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	wish "github.com/charmbracelet/wish"
	wishbubble "github.com/charmbracelet/wish/bubbletea"
	wishlog "github.com/charmbracelet/wish/logging"
	"github.com/pavandhadge/ssh_portfolio_opensource/internal/tui"
)

func envOrDefault(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func ensureHostKey(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}
	key := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	return os.WriteFile(path, key, 0o600)
}

func teaHandler(_ ssh.Session) (tea.Model, []tea.ProgramOption) {
	return tui.NewModel(), []tea.ProgramOption{tea.WithMouseCellMotion()}
}

func newWishServer(addr, hostKeyPath string) (*ssh.Server, error) {
	if err := ensureHostKey(hostKeyPath); err != nil {
		return nil, fmt.Errorf("prepare host key: %w", err)
	}

	server, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(hostKeyPath),
		wish.WithIdleTimeout(5*time.Minute),
		wish.WithMaxTimeout(10*time.Minute),
		wish.WithMiddleware(
			wishlog.Middleware(),
			wishbubble.Middleware(teaHandler),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create SSH server: %w", err)
	}
	return server, nil
}

func Run() error {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			return fmt.Errorf("debug log setup: %w", err)
		}
		defer f.Close()
	}

	host := envOrDefault("SSH_PORTFOLIO_HOST", "127.0.0.1")
	port := envOrDefault("SSH_PORTFOLIO_PORT", "2222")
	addr := net.JoinHostPort(host, port)
	hostKeyPath := envOrDefault("SSH_PORTFOLIO_HOST_KEY", ".wish/ssh_portfolio_ed25519")

	server, err := newWishServer(addr, hostKeyPath)
	if err != nil {
		return err
	}

	go func() {
		log.Printf("ssh portfolio listening on ssh://%s (OpenSSH on :22 unchanged)", addr)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("ssh server stopped: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return server.Shutdown(ctx)
}
