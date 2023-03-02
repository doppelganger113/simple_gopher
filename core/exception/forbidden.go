package exception

const defaultReason = "Forbidden"

type Forbidden struct {
	Reason string
}

func (f Forbidden) Error() string {
	if f.Reason == "" {
		return defaultReason
	}
	return f.Reason
}
