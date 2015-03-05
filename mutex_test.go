package main

import (
	"net"
	"testing"
)

func TestNoAddresses(t *testing.T) {
	ui := &SilentLog{}
	addrs := []net.Addr{}
	_, err := NewMutex("keyname", addrs, ui)
	expected := "redis: addrs is empty"
	if err.Error() != expected {
		t.Errorf("got '%s', expected '%v'", err.Error(), expected)
	}
}

func TestCantConnect(t *testing.T) {
	ui := &SilentLog{}
	_, err := NewMutex("keyname", []net.Addr{
		&net.TCPAddr{Port: 6379, IP: net.ParseIP("10.0.0.0")},
	}, ui)
	expected := "Failed to connect to any redis server"
	if err == nil {
		t.Errorf("no connectableredis server should have thrown an error")
		return
	}

	if err.Error() != expected {
		t.Errorf("got '%s', expected '%v'", err.Error(), expected)
	}
}

func TestCanConnect(t *testing.T) {
	ui := &SilentLog{}
	m, err := NewMutex("keyname", []net.Addr{
		&net.TCPAddr{Port: 6379, IP: net.ParseIP("127.0.0.1")},
	}, ui)

	if err != nil {
		t.Errorf("should have been able to connect to 127.0.0.1:6370")
		return
	}

	expected := 1
	got := len(m.Nodes)
	if got != expected {
		t.Errorf("got '%s', expected '%s'", got, expected)
		return
	}

	expected = 1
	got = m.Quorum
	if got != expected {
		t.Errorf("got '%s', expected '%s'", got, expected)
		return
	}

	expectedName := "keyname"
	gotName := m.Name
	if got != expected {
		t.Errorf("got '%s', expected '%s'", expectedName, gotName)
		return
	}
}

func TestCanLock(t *testing.T) {
	ui := &SilentLog{}
	m, _ := NewMutex("keyname", []net.Addr{
		&net.TCPAddr{Port: 6379, IP: net.ParseIP("127.0.0.1")},
	}, ui)

	err := m.Lock()
	defer m.Unlock()

	if err != nil {
		t.Errorf("should have aquired lock")
		return
	}

	if m.Value() == "" {
		t.Errorf("m.value should not be empty")
		return
	}

	if m.Until().IsZero() {
		t.Errorf("m.until should not be zero")
		return
	}
}
