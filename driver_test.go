package main

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestListClosesControllerConnections(t *testing.T) {
	const calls = 8
	var closedCount int32
	closed := make(chan struct{})
	controller := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, "[]")
	}))
	controller.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateClosed && atomic.AddInt32(&closedCount, 1) == calls {
			close(closed)
		}
	}
	controller.Start()
	defer controller.Close()

	dir := t.TempDir()
	config := filepath.Join(dir, "docker-volume.conf")
	if err := os.WriteFile(config, []byte("[global]\ncontrollers = "+controller.URL+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	driver := NewLinstorDriver(config, "test-node", dir)
	for i := 0; i < calls; i++ {
		if _, err := driver.List(); err != nil {
			t.Fatal(err)
		}
	}

	timeout := time.NewTimer(2 * time.Second)
	defer timeout.Stop()
	select {
	case <-closed:
	case <-timeout.C:
		t.Fatalf("closed only %d of %d controller connections", atomic.LoadInt32(&closedCount), calls)
	}
}
