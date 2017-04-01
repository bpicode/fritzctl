package fritz

// MessageLog is the top-level wrapper for the FRITZ!Box answer upon
// a query for log events.
type MessageLog struct {
	Messages []Message `json:"mq_log"`
}

// Message corresponds to a single log message.
type Message []string
