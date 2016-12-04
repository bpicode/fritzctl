package fritz

// MessageLog is the top-level wrapper for the FRITZ!Box answer upon
// a query for log events.
type MessageLog struct {
	Messages []Message `json:"mq_log"`
}

// Message corresponds to a single log message.
type Message struct {
	Text string `json:"_node"`
	Row2 string `json:"row_2"`
	Row3 string `json:"row_3"`
}
