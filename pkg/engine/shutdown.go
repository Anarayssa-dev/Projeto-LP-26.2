package engine

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sort"
	"syscall"
)

// MultiPingFunc representa a implementação do motor paralelo.
// A função deve respeitar o contexto e fechar resultsChan somente depois
// que todas as goroutines de ping terminarem.
type MultiPingFunc func(ctx context.Context, targets []net.IP, resultsChan chan<- PingResult)

// RunWithGracefulShutdown executa o motor e escuta Ctrl+C (SIGINT).
// Ao receber o sinal, cancela o contexto, continua consumindo os resultados
// já produzidos e imprime as métricas quando o canal for fechado.
func RunWithGracefulShutdown(targets []net.IP, multiPing MultiPingFunc, output io.Writer) error {
	if multiPing == nil {
		return fmt.Errorf("MultiPing não pode ser nil")
	}
	if output == nil {
		output = os.Stdout
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	return runWithSignals(targets, multiPing, output, signals)
}

// runWithSignals fica separado para permitir testes sem enviar sinais reais ao processo.
func runWithSignals(targets []net.IP, multiPing MultiPingFunc, output io.Writer, signals <-chan os.Signal) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultsChan := make(chan PingResult)
	done := make(chan struct{})
	go func() {
		defer close(done)
		multiPing(ctx, targets, resultsChan)
	}()

	stats := make(map[string]*HostStats)
	interrupted := false
	for {
		select {
		case result, ok := <-resultsChan:
			if !ok {
				printReport(output, stats)
				<-done
				return nil
			}
			host := stats[result.Target]
			if host == nil {
				host = &HostStats{IP: result.Target}
				stats[result.Target] = host
			}
			host.AddResult(result)
		case <-signals:
			if !interrupted {
				interrupted = true
				cancel()
				fmt.Fprintln(output, "\nInterrupção recebida. Encerrando os pings...")
			}
			// Continua consumindo até MultiPing fechar o canal, evitando bloquear envios.
		}
	}
}

func printReport(output io.Writer, stats map[string]*HostStats) {
	fmt.Fprintln(output, "\n===== RELATÓRIO DE MÉTRICAS =====")
	if len(stats) == 0 {
		fmt.Fprintln(output, "Nenhum resultado de ping foi recebido.")
		return
	}

	hosts := make([]string, 0, len(stats))
	for host := range stats {
		hosts = append(hosts, host)
	}
	sort.Strings(hosts)
	for _, host := range hosts {
		s := stats[host]
		fmt.Fprintf(output, "%s | Pacotes: %d | RTT mín: %s | RTT máx: %s | RTT médio: %s | Perda: %.2f%%\n",
			host, len(s.History), s.RTTMin(), s.RTTMax(), s.RTTAvg(), s.PacketLoss())
	}
}
