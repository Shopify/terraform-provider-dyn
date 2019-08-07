package dyn

// ZoneOption is a basic interface for passing optional parameters for a zone
type ZoneOption func(*zoneCreateRequest)

// SerialStyle provides a serialStyle option
func SerialStyle(serialStyle string) ZoneOption {
	return func(r *zoneCreateRequest) {
		r.SerialStyle = serialStyle
	}
}
