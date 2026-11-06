
package engine

import "time"

type PingResult struct {
    Target    string
    RTT       time.Duration
    Success   bool
    Timestamp time.Time
}

type MetricsStorage interface {
    Save(result PingResult) error
    GetAll(target string) ([]PingResult, error)
    Close() error
}

// Armazena o histórico de disparos de um IP.
type HostStats struct {
    IP      string
    History []PingResult
}

// Adiciona um resultado ao histórico.
func (h *HostStats) AddResult(result PingResult) {
    if h.IP == "" {
        h.IP = result.Target
    }

    h.History = append(h.History, result)
}

// Calcula o RTT mínimo.
func (h *HostStats) RTTMin() time.Duration {
    var min time.Duration
    found := false

    for _, result := range h.History {
        if result.Success {
            if !found || result.RTT < min {
                min = result.RTT
                found = true
            }
        }
    }

    return min
}

// Calcula o RTT máximo.
func (h *HostStats) RTTMax() time.Duration {
    var max time.Duration
    found := false

    for _, result := range h.History {
        if result.Success {
            if !found || result.RTT > max {
                max = result.RTT
                found = true
            }
        }
    }

    return max
}

// Calcula o RTT médio.
func (h *HostStats) RTTAvg() time.Duration {
    var total time.Duration
    var count int

    for _, result := range h.History {
        if result.Success {
            total += result.RTT
            count++
        }
    }

    if count == 0 {
        return 0
    }

    return total / time.Duration(count)
}

// Calcula a porcentagem de perda de pacotes.
func (h *HostStats) PacketLoss() float64 {
    total := len(h.History)

    if total == 0 {
        return 0
    }

    failures := 0

    for _, result := range h.History {
        if !result.Success {
            failures++
        }
    }

    return float64(failures) / float64(total) * 100
}
