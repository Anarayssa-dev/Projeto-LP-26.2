package main

import(
	"fmt"
	"net"
	"errors"
)

func ResolveTargets(hosts []string) (map[string][]net.IP, error){

	//Slice que vai guardar os ips resultantes
	ips := map[string][]net.IP{}
	failedTargets := []string{}
	for _,target := range hosts {
		//Se o alvo for um ip direto
		ip := net.ParseIP(target)

		if ip != nil{
			ips[target] = []net.IP{ip}
		}

		//Se o alvo for um domínio
		domainIPs, err := net.LookupIP(target)//net.LookupIP retorna um slice com os ips do domínio

		if err == nil{
			ips[target] = domainIPs
		}

		//Se ambos falharem, o alvo não foi resolvido
		if ip == nil && err != nil{
			failedTargets = append(failedTargets, target)
		}
	}

	if len(failedTargets)!=0{
		return ips, errors.New("Argumento(s) inválido(s): " + fmt.Sprint(failedTargets))
	}
	
	return ips, nil
}