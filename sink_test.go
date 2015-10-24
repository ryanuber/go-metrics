package metrics

import (
	"reflect"
	"testing"
)

type MockSink struct {
	keys [][]string
	vals []float32
}

func (m *MockSink) SetGauge(key []string, val float32) {
	m.keys = append(m.keys, key)
	m.vals = append(m.vals, val)
}
func (m *MockSink) EmitKey(key []string, val float32) {
	m.keys = append(m.keys, key)
	m.vals = append(m.vals, val)
}
func (m *MockSink) IncrCounter(key []string, val float32) {
	m.keys = append(m.keys, key)
	m.vals = append(m.vals, val)
}
func (m *MockSink) AddSample(key []string, val float32) {
	m.keys = append(m.keys, key)
	m.vals = append(m.vals, val)
}

func TestFanoutSink_Gauge(t *testing.T) {
	m1 := &MockSink{}
	m2 := &MockSink{}
	fh := &FanoutSink{m1, m2}

	k := []string{"test"}
	v := float32(42.0)
	fh.SetGauge(k, v)

	if !reflect.DeepEqual(m1.keys[0], k) {
		t.Fatalf("key not equal")
	}
	if !reflect.DeepEqual(m2.keys[0], k) {
		t.Fatalf("key not equal")
	}
	if !reflect.DeepEqual(m1.vals[0], v) {
		t.Fatalf("val not equal")
	}
	if !reflect.DeepEqual(m2.vals[0], v) {
		t.Fatalf("val not equal")
	}
}

func TestFanoutSink_Key(t *testing.T) {
	m1 := &MockSink{}
	m2 := &MockSink{}
	fh := &FanoutSink{m1, m2}

	k := []string{"test"}
	v := float32(42.0)
	fh.EmitKey(k, v)

	if !reflect.DeepEqual(m1.keys[0], k) {
		t.Fatalf("key not equal")
	}
	if !reflect.DeepEqual(m2.keys[0], k) {
		t.Fatalf("key not equal")
	}
	if !reflect.DeepEqual(m1.vals[0], v) {
		t.Fatalf("val not equal")
	}
	if !reflect.DeepEqual(m2.vals[0], v) {
		t.Fatalf("val not equal")
	}
}

func TestFanoutSink_Counter(t *testing.T) {
	m1 := &MockSink{}
	m2 := &MockSink{}
	fh := &FanoutSink{m1, m2}

	k := []string{"test"}
	v := float32(42.0)
	fh.IncrCounter(k, v)

	if !reflect.DeepEqual(m1.keys[0], k) {
		t.Fatalf("key not equal")
	}
	if !reflect.DeepEqual(m2.keys[0], k) {
		t.Fatalf("key not equal")
	}
	if !reflect.DeepEqual(m1.vals[0], v) {
		t.Fatalf("val not equal")
	}
	if !reflect.DeepEqual(m2.vals[0], v) {
		t.Fatalf("val not equal")
	}
}

func TestFanoutSink_Sample(t *testing.T) {
	m1 := &MockSink{}
	m2 := &MockSink{}
	fh := &FanoutSink{m1, m2}

	k := []string{"test"}
	v := float32(42.0)
	fh.AddSample(k, v)

	if !reflect.DeepEqual(m1.keys[0], k) {
		t.Fatalf("key not equal")
	}
	if !reflect.DeepEqual(m2.keys[0], k) {
		t.Fatalf("key not equal")
	}
	if !reflect.DeepEqual(m1.vals[0], v) {
		t.Fatalf("val not equal")
	}
	if !reflect.DeepEqual(m2.vals[0], v) {
		t.Fatalf("val not equal")
	}
}

func TestFilterSink(t *testing.T) {
	// Create the filters
	filt := []FilterFunc{
		func(key []string, val float32) bool { return key[0] == "baz" },
		func(key []string, val float32) bool { return val < 1.0 },
	}
	fs := &FilterSink{Filters: filt}

	funcs := []func(key []string, val float32){
		fs.SetGauge,
		fs.EmitKey,
		fs.IncrCounter,
		fs.AddSample,
	}

	for _, fn := range funcs {
		m1 := &MockSink{}
		fs.Sink = m1

		// Trigger some metrics
		fn([]string{"foo"}, 2.9) // Passes
		fn([]string{"bar"}, 1.3) // Passes
		fn([]string{"baz"}, 9.6) // Filters key
		fn([]string{"zip"}, 0.6) // Filters val

		expectKeys := [][]string{[]string{"foo"}, []string{"bar"}}
		if !reflect.DeepEqual(m1.keys, expectKeys) {
			t.Fatalf("bad keys: expect %v, got: %v", expectKeys, m1.keys)
		}

		expectVals := []float32{2.9, 1.3}
		if !reflect.DeepEqual(m1.vals, expectVals) {
			t.Fatalf("bad vals: expect %v, got: %v", expectVals, m1.vals)
		}
	}
}
