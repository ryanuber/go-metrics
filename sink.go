package metrics

// The MetricSink interface is used to transmit metrics information
// to an external system
type MetricSink interface {
	// A Gauge should retain the last value it is set to
	SetGauge(key []string, val float32)

	// Should emit a Key/Value pair for each call
	EmitKey(key []string, val float32)

	// Counters should accumulate values
	IncrCounter(key []string, val float32)

	// Samples are for timing information, where quantiles are used
	AddSample(key []string, val float32)
}

// BlackholeSink is used to just blackhole messages
type BlackholeSink struct{}

func (*BlackholeSink) SetGauge(key []string, val float32)    {}
func (*BlackholeSink) EmitKey(key []string, val float32)     {}
func (*BlackholeSink) IncrCounter(key []string, val float32) {}
func (*BlackholeSink) AddSample(key []string, val float32)   {}

// FanoutSink is used to sink to fanout values to multiple sinks
type FanoutSink []MetricSink

func (fh FanoutSink) SetGauge(key []string, val float32) {
	for _, s := range fh {
		s.SetGauge(key, val)
	}
}

func (fh FanoutSink) EmitKey(key []string, val float32) {
	for _, s := range fh {
		s.EmitKey(key, val)
	}
}

func (fh FanoutSink) IncrCounter(key []string, val float32) {
	for _, s := range fh {
		s.IncrCounter(key, val)
	}
}

func (fh FanoutSink) AddSample(key []string, val float32) {
	for _, s := range fh {
		s.AddSample(key, val)
	}
}

// FilterFunc is a function interface used to determine if the given metrics
// key or value should cause the metrics layer to filter the message.
type FilterFunc func(key []string, val float32) bool

// FilterSink is an implementation of MetricSink which allows pre-filtering
// of submitted metrics. If the filters pass, then the underlink MetricSink
// is handed the values normally.
type FilterSink struct {
	Sink    MetricSink
	Filters []FilterFunc
}

func (s *FilterSink) SetGauge(key []string, val float32) {
	if !s.filter(key, val) {
		s.Sink.SetGauge(key, val)
	}
}

func (s *FilterSink) EmitKey(key []string, val float32) {
	if !s.filter(key, val) {
		s.Sink.EmitKey(key, val)
	}
}

func (s *FilterSink) IncrCounter(key []string, val float32) {
	if !s.filter(key, val) {
		s.Sink.IncrCounter(key, val)
	}
}

func (s *FilterSink) AddSample(key []string, val float32) {
	if !s.filter(key, val) {
		s.Sink.AddSample(key, val)
	}
}

// filter is used to iterate our filters and check if they would filter
// the given metrics value.
func (s *FilterSink) filter(key []string, val float32) bool {
	for _, fn := range s.Filters {
		if fn(key, val) {
			return true
		}
	}
	return false
}
