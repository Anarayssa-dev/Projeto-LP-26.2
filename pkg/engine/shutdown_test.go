package engine

import (
	"bytes"
	"context"
	"net"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestRunWithSignalsPrintsAccumulatedMetricsAfterInterrupt(t *testing.T) {
	var output bytes.Buffer
	signals := make(chan os.Signal, 1)
	started := make(chan struct{})

	multiPing := func(ctx context.Context, _ []net.IP, results chan<- PingResult) {
		defer close(results)
		results <- PingResult{Target: "192.0.2.1", RTT: 10 * time.Millisecond, Success: true}
		close(started)
		<-ctx.Done()
		results <- PingResult{Target: "192.0.2.1", Success: false}
	}

	done := make(chan error, 1)
	go func() {
		done <- runWithSignals(nil, multiPing, &output, signals)
	}()
	<-started
	signals <- syscall.SIGINT

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("runWithSignals() erro: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("execução não terminou após o cancelamento")
	}

	got := output.String()
	for _, want := range []string{"Interrupção recebida", "RELATÓRIO DE MÉTRICAS", "192.0.2.1", "Pacotes: 2", "Perda: 50.00%"} {
		if !strings.Contains(got, want) {
			t.Errorf("saída não contém %q:\n%s", want, got)
		}
	}
}

func TestRunWithGracefulShutdownRejectsNilMultiPing(t *testing.T) {
	if err := RunWithGracefulShutdown(nil, nil, &bytes.Buffer{}); err == nil {
		t.Fatal("esperava erro quando MultiPing é nil")
	}
}
