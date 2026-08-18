package app

import "testing"

func TestNewApp(t *testing.T) {
	if app := NewApp(); app == nil {
		t.Errorf("NewApp() expected not to be nil")
	}
}

func TestCmd(t *testing.T) {
	if err := testCmd.RunE(nil, nil); err == nil {
		t.Errorf("expected error but got nil instead")
	}
}
