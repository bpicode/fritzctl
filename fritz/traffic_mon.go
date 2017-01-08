package fritz

// TrafficMonitoringData holds the data for the up- and
// downstream traffic reported by the FRITZ!Box.
type TrafficMonitoringData struct {
	DownstreamInternetBitsPerSecond      []float64 `json:"ds_current_bps"`
	DownStreamMediaBitsPerSecond         []float64 `json:"mc_current_bps"`
	UpstreamRealtimeBitsPerSecond        []float64 `json:"prio_realtime_bps"`
	UpstreamHighPriorityBitsPerSecond    []float64 `json:"prio_high_bps"`
	UpstreamLowPriorityBitsPerSecond     []float64 `json:"prio_low_bps"`
	UpstreamDefaultPriorityBitsPerSecond []float64 `json:"prio_default_bps"`
}
