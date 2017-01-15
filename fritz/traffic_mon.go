package fritz

// TrafficMonitoringData holds the data for the up- and downstream traffic reported by the FRITZ!Box.
// codebeat:disable[TOO_MANY_IVARS]
type TrafficMonitoringData struct {
	DownstreamInternet      []float64 `json:"ds_current_bps"`
	DownStreamMedia         []float64 `json:"mc_current_bps"`
	UpstreamRealtime        []float64 `json:"prio_realtime_bps"`
	UpstreamHighPriority    []float64 `json:"prio_high_bps"`
	UpstreamDefaultPriority []float64 `json:"prio_default_bps"`
	UpstreamLowPriority     []float64 `json:"prio_low_bps"`
}

// codebeat:enable[TOO_MANY_IVARS]

// BitsPerSecond returns a TrafficMonitoringData with metrics in units of bits/second.
func (d TrafficMonitoringData) BitsPerSecond() TrafficMonitoringData {
	return TrafficMonitoringData{
		DownstreamInternet:      d.DownstreamInternet,
		DownStreamMedia:         d.DownStreamMedia,
		UpstreamRealtime:        d.UpstreamRealtime,
		UpstreamHighPriority:    d.UpstreamHighPriority,
		UpstreamDefaultPriority: d.UpstreamDefaultPriority,
		UpstreamLowPriority:     d.UpstreamLowPriority,
	}
}

// KiloBitsPerSecond returns a TrafficMonitoringData with metrics in units of kbits/second.
func (d TrafficMonitoringData) KiloBitsPerSecond() TrafficMonitoringData {
	return TrafficMonitoringData{
		DownstreamInternet:      div(d.DownstreamInternet, 1024),
		DownStreamMedia:         div(d.DownStreamMedia, 1024),
		UpstreamRealtime:        div(d.UpstreamRealtime, 1024),
		UpstreamHighPriority:    div(d.UpstreamHighPriority, 1024),
		UpstreamDefaultPriority: div(d.UpstreamDefaultPriority, 1024),
		UpstreamLowPriority:     div(d.UpstreamLowPriority, 1024),
	}
}

func div(xs []float64, d float64) []float64 {
	ys := make([]float64, len(xs))
	for i, x := range xs {
		ys[i] = x / d
	}
	return ys
}
