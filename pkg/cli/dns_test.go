package main

import (
	"errors"
	"net"
	"reflect"
	"testing"
)

func withLookupIP(t *testing.T, lookup func(string) ([]net.IP, error)) {
	t.Helper()
	original := lookupIP
	lookupIP = lookup
	t.Cleanup(func() { lookupIP = original })
}

func TestResolveTargetsAcceptsIPv4(t *testing.T) {
	withLookupIP(t, func(string) ([]net.IP, error) {
		t.Fatal("não deveria consultar DNS para um endereço IPv4")
		return nil, nil
	})

	got, err := ResolveTargets([]string{"192.0.2.10"})
	if err != nil {
		t.Fatalf("ResolveTargets() erro inesperado: %v", err)
	}
	want := net.ParseIP("192.0.2.10")
	if ips := got["192.0.2.10"]; len(ips) != 1 || !ips[0].Equal(want) {
		t.Fatalf("resultado IPv4 inesperado: %#v", got)
	}
}

func TestResolveTargetsAcceptsIPv6(t *testing.T) {
	withLookupIP(t, func(string) ([]net.IP, error) {
		t.Fatal("não deveria consultar DNS para um endereço IPv6")
		return nil, nil
	})

	got, err := ResolveTargets([]string{"2001:db8::1"})
	if err != nil {
		t.Fatalf("ResolveTargets() erro inesperado: %v", err)
	}
	want := net.ParseIP("2001:db8::1")
	if ips := got["2001:db8::1"]; len(ips) != 1 || !ips[0].Equal(want) {
		t.Fatalf("resultado IPv6 inesperado: %#v", got)
	}
}

func TestResolveTargetsResolvesDomain(t *testing.T) {
	expected := []net.IP{net.ParseIP("192.0.2.30"), net.ParseIP("2001:db8::30")}
	withLookupIP(t, func(host string) ([]net.IP, error) {
		if host != "exemplo.test" {
			t.Fatalf("domínio consultado inesperado: %q", host)
		}
		return expected, nil
	})

	got, err := ResolveTargets([]string{"exemplo.test"})
	if err != nil {
		t.Fatalf("ResolveTargets() erro inesperado: %v", err)
	}
	if !reflect.DeepEqual(got["exemplo.test"], expected) {
		t.Fatalf("IPs resolvidos = %#v; esperado %#v", got["exemplo.test"], expected)
	}
}

func TestResolveTargetsRejectsInvalidDomain(t *testing.T) {
	withLookupIP(t, func(string) ([]net.IP, error) {
		return nil, errors.New("DNS lookup falhou")
	})

	got, err := ResolveTargets([]string{"nao-existe.test"})
	if err == nil {
		t.Fatal("esperava erro para domínio não resolvido")
	}
	if len(got) != 0 {
		t.Fatalf("não esperava IPs resolvidos: %#v", got)
	}
}

func TestResolveTargetsKeepsValidTargetsWhenAnotherFails(t *testing.T) {
	withLookupIP(t, func(host string) ([]net.IP, error) {
		if host == "invalido.test" {
			return nil, errors.New("DNS lookup falhou")
		}
		return []net.IP{net.ParseIP("192.0.2.20")}, nil
	})

	got, err := ResolveTargets([]string{"192.0.2.20", "invalido.test"})
	if err == nil {
		t.Fatal("esperava erro por causa do domínio inválido")
	}
	if ips := got["192.0.2.20"]; len(ips) != 1 || !ips[0].Equal(net.ParseIP("192.0.2.20")) {
		t.Fatalf("o IP válido deveria permanecer no resultado: %#v", got)
	}
}
