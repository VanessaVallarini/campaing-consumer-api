package health

type MonitoredStatus string

const (
	StatusUp      MonitoredStatus = "UP"
	StatusDown    MonitoredStatus = "DOWN"
	StatusUnknown MonitoredStatus = "UNKNOWN"
)

// Rendered output of healthcheck calls
type healthCheckData struct {
	Status     MonitoredStatus            `json:"status"`
	Components map[string]MonitoredStatus `json:"components"`
}
