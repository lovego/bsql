package scan

import (
	"reflect"
	"testing"
)

func TestJsonScannerNilValue(t *testing.T) {
	var m map[string]int
	js := jsonScanner{destAddr: reflect.ValueOf(&m)}
	if err := js.Scan(nil); err != nil {
		t.Fatal(err)
	}
	if m != nil {
		t.Errorf("unexpected: %v", reflect.ValueOf(m))
	}

	m = map[string]int{"key": 1}
	js = jsonScanner{destAddr: reflect.ValueOf(&m)}
	if err := js.Scan(nil); err != nil {
		t.Fatal(err)
	}
	if m != nil {
		t.Errorf("unexpected: %v", m)
	}
}
