package bsql

import (
	"testing"
)

func TestJsonScannerNilValue(t *testing.T) {
	var m map[string]int
	js := jsonScanner{dest: &m}
	if err := js.Scan(nil); err != nil {
		t.Fatal(err)
	}
	if m != nil {
		t.Errorf("unexpected: %v", m)
	}

	m = map[string]int{"key": 1}
	js = jsonScanner{dest: &m}
	if err := js.Scan(nil); err != nil {
		t.Fatal(err)
	}
	if m != nil {
		t.Errorf("unexpected: %v", m)
	}
}
