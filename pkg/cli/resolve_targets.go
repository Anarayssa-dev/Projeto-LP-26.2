package main

import (
	"errors"
	"fmt"
	"net"
)

// lookupIP permite testar a resolução de nomes sem depender de um servidor DNS.
var lookupIP = net.LookupIP

func ResolveTargets(hosts []string) (map[string][]net.IP, error) {
	ips := make(map[string][]net.IP)
	var failedTargets []string

	for _, target := range hosts {
		// Se o alvo já for um IP, não é necessário consultar o DNS.
		if ip := net.ParseIP(target); ip != nil {
			ips[target] = []net.IP{ip}
			continue
		}

		// Se o alvo for um domínio, tenta resolver seu endereço.
		domainIPs, err := lookupIP(target)
		if err == nil && len(domainIPs) > 0 {
			ips[target] = domainIPs
			continue
		}

		failedTargets = append(failedTargets, target)
	}

	if len(failedTargets) > 0 {
		return ips, errors.New("Argumento(s) inválido(s): " + fmt.Sprint(failedTargets))
	}

	return ips, nil
}
