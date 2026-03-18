package server

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	cryptossh "golang.org/x/crypto/ssh"
)

func envIntBounded(key string, def, minV, maxV int) int {
	v := def
	if raw := os.Getenv(key); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil {
			v = parsed
		}
	}
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

func waitTCPReady(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("server not ready at %s within %s", addr, timeout)
}

type loadMetrics struct {
	Total      int
	Parallel   int
	OK         int
	Fail       int
	SuccessPct int
	Duration   time.Duration
	Throughput float64
	P50        time.Duration
	P95        time.Duration
	P99        time.Duration
}

func runLoadRound(addr string, total, parallel int, timeout time.Duration) loadMetrics {
	cfg := &cryptossh.ClientConfig{
		User:            "loadtest",
		Auth:            []cryptossh.AuthMethod{},
		HostKeyCallback: cryptossh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}

	sem := make(chan struct{}, parallel)
	var okCount int64
	var failCount int64
	latencies := make([]time.Duration, total)
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			connStart := time.Now()
			client, err := cryptossh.Dial("tcp", addr, cfg)
			if err != nil {
				atomic.AddInt64(&failCount, 1)
				return
			}
			defer client.Close()

			sess, err := client.NewSession()
			if err != nil {
				atomic.AddInt64(&failCount, 1)
				return
			}
			_ = sess.Close()
			latencies[idx] = time.Since(connStart)
			atomic.AddInt64(&okCount, 1)
		}(i)
	}
	wg.Wait()
	dur := time.Since(start)

	ok := int(atomic.LoadInt64(&okCount))
	fail := int(atomic.LoadInt64(&failCount))
	successPct := 0
	if total > 0 {
		successPct = (ok * 100) / total
	}

	lat := make([]time.Duration, 0, ok)
	for _, d := range latencies {
		if d > 0 {
			lat = append(lat, d)
		}
	}
	sort.Slice(lat, func(i, j int) bool { return lat[i] < lat[j] })
	percentile := func(p float64) time.Duration {
		if len(lat) == 0 {
			return 0
		}
		pos := int(math.Ceil((p/100.0)*float64(len(lat)))) - 1
		if pos < 0 {
			pos = 0
		}
		if pos >= len(lat) {
			pos = len(lat) - 1
		}
		return lat[pos]
	}
	throughput := 0.0
	if dur > 0 {
		throughput = float64(ok) / dur.Seconds()
	}

	return loadMetrics{
		Total:      total,
		Parallel:   parallel,
		OK:         ok,
		Fail:       fail,
		SuccessPct: successPct,
		Duration:   dur,
		Throughput: throughput,
		P50:        percentile(50),
		P95:        percentile(95),
		P99:        percentile(99),
	}
}

func startTestServer(t *testing.T) (addr string, cleanup func()) {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen for test port: %v", err)
	}
	addr = l.Addr().String()
	_ = l.Close()

	keyPath := filepath.Join(t.TempDir(), "wish_test_ed25519")
	srv, err := newWishServer(addr, keyPath)
	if err != nil {
		t.Fatalf("create test wish server: %v", err)
	}

	serverErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	if err := waitTCPReady(addr, 4*time.Second); err != nil {
		t.Fatal(err)
	}

	cleanup = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)

		select {
		case err := <-serverErr:
			if err != nil && !errors.Is(err, net.ErrClosed) && !strings.Contains(err.Error(), "Server closed") {
				t.Fatalf("server exited with error: %v", err)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for server to stop")
		}
	}
	return addr, cleanup
}

func TestWishConcurrentConnectionsSafeBurst(t *testing.T) {
	t.Parallel()

	total := envIntBounded("SSH_PORTFOLIO_LOAD_TOTAL", 80, 10, 300)
	parallel := envIntBounded("SSH_PORTFOLIO_LOAD_PARALLEL", 24, 2, 64)
	successThresholdPct := envIntBounded("SSH_PORTFOLIO_LOAD_SUCCESS_PCT", 85, 50, 100)
	timeoutMs := envIntBounded("SSH_PORTFOLIO_LOAD_TIMEOUT_MS", 2000, 500, 5000)

	addr, cleanup := startTestServer(t)
	defer cleanup()

	m := runLoadRound(addr, total, parallel, time.Duration(timeoutMs)*time.Millisecond)

	aliveConn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
	if err != nil {
		t.Fatalf("server stopped responding under load: %v", err)
	}
	_ = aliveConn.Close()

	if m.SuccessPct < successThresholdPct {
		t.Fatalf("load test failed: ok=%d fail=%d total=%d success=%d%% threshold=%d%%", m.OK, m.Fail, m.Total, m.SuccessPct, successThresholdPct)
	}

	t.Logf("concurrency result: ok=%d fail=%d total=%d success=%d%% parallel=%d dur=%s throughput=%.2fconn/s p50=%s p95=%s p99=%s",
		m.OK, m.Fail, m.Total, m.SuccessPct, m.Parallel, m.Duration, m.Throughput, m.P50, m.P95, m.P99)
}

func TestWishConcurrentConnectionsRampSafe(t *testing.T) {
	startParallel := envIntBounded("SSH_PORTFOLIO_RAMP_START_PAR", 12, 2, 64)
	stepParallel := envIntBounded("SSH_PORTFOLIO_RAMP_STEP_PAR", 8, 1, 32)
	maxParallel := envIntBounded("SSH_PORTFOLIO_RAMP_MAX_PAR", 64, startParallel, 128)
	multiplier := envIntBounded("SSH_PORTFOLIO_RAMP_TOTAL_MULT", 4, 2, 8)
	maxTotal := envIntBounded("SSH_PORTFOLIO_RAMP_MAX_TOTAL", 512, 50, 2000)
	successThresholdPct := envIntBounded("SSH_PORTFOLIO_LOAD_SUCCESS_PCT", 85, 50, 100)
	timeoutMs := envIntBounded("SSH_PORTFOLIO_LOAD_TIMEOUT_MS", 2000, 500, 5000)

	addr, cleanup := startTestServer(t)
	defer cleanup()

	t.Logf("ramp config: start_par=%d step=%d max_par=%d total_mult=%d max_total=%d threshold=%d%% timeout_ms=%d",
		startParallel, stepParallel, maxParallel, multiplier, maxTotal, successThresholdPct, timeoutMs)

	best := loadMetrics{}
	for par := startParallel; par <= maxParallel; par += stepParallel {
		total := par * multiplier
		if total > maxTotal {
			total = maxTotal
		}
		m := runLoadRound(addr, total, par, time.Duration(timeoutMs)*time.Millisecond)
		t.Logf("round: par=%d total=%d ok=%d fail=%d success=%d%% dur=%s throughput=%.2fconn/s p50=%s p95=%s p99=%s",
			m.Parallel, m.Total, m.OK, m.Fail, m.SuccessPct, m.Duration, m.Throughput, m.P50, m.P95, m.P99)

		if m.SuccessPct >= successThresholdPct {
			best = m
		} else {
			t.Logf("ramp stop: success %d%% fell below threshold %d%% at par=%d", m.SuccessPct, successThresholdPct, par)
			break
		}
		time.Sleep(150 * time.Millisecond)
	}

	if best.Total == 0 {
		t.Fatalf("no safe concurrency level met threshold=%d%%", successThresholdPct)
	}
	t.Logf("safe limit: parallel=%d total=%d success=%d%% throughput=%.2fconn/s p95=%s p99=%s",
		best.Parallel, best.Total, best.SuccessPct, best.Throughput, best.P95, best.P99)
}
