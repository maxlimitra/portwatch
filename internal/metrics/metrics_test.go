package metrics

import (
	"bytes"
	"strings"
	"testing"
)

func TestCounterIncAndValue(t *testing.T) {
	c := &Counter{}
	if c.Value() != 0 {
		t.Fatalf("expected 0, got %d", c.Value())
	}
	c.Inc()
	c.Inc()
	c.Add(3)
	if c.Value() != 5 {
		t.Fatalf("expected 5, got %d", c.Value())
	}
}

func TestGaugeSetAndValue(t *testing.T) {
	g := &Gauge{}
	g.Set(42)
	if g.Value() != 42 {
		t.Fatalf("expected 42, got %d", g.Value())
	}
	g.Set(-7)
	if g.Value() != -7 {
		t.Fatalf("expected -7, got %d", g.Value())
	}
}

func TestRegistryCounterSameInstance(t *testing.T) {
	r := NewRegistry()
	c1 := r.Counter("scans")
	c1.Inc()
	c2 := r.Counter("scans")
	if c2.Value() != 1 {
		t.Fatalf("expected same instance; got value %d", c2.Value())
	}
}

func TestRegistryGaugeSameInstance(t *testing.T) {
	r := NewRegistry()
	g1 := r.Gauge("ports_tracked")
	g1.Set(10)
	g2 := r.Gauge("ports_tracked")
	if g2.Value() != 10 {
		t.Fatalf("expected same instance; got value %d", g2.Value())
	}
}

func TestWriteToContainsAllMetrics(t *testing.T) {
	r := NewRegistry()
	r.Counter("alerts_fired").Add(3)
	r.Counter("scan_cycles").Inc()
	r.Gauge("ports_tracked").Set(7)

	var buf bytes.Buffer
	r.WriteTo(&buf)
	out := buf.String()

	for _, want := range []string{
		"counter alerts_fired 3",
		"counter scan_cycles 1",
		"gauge   ports_tracked 7",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\ngot:\n%s", want, out)
		}
	}
}

func TestWriteToSortedOutput(t *testing.T) {
	r := NewRegistry()
	r.Counter("z_last").Inc()
	r.Counter("a_first").Inc()

	var buf bytes.Buffer
	r.WriteTo(&buf)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")

	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "a_first") {
		t.Errorf("expected a_first first, got: %s", lines[0])
	}
	if !strings.Contains(lines[1], "z_last") {
		t.Errorf("expected z_last second, got: %s", lines[1])
	}
}
