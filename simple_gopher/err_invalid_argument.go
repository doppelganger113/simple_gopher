package simple_gopher

type InvalidArgument struct {
	Reason string
}

func (ia InvalidArgument) Error() string {
	return ia.Reason
}
