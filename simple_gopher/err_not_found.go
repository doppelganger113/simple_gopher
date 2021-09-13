package simple_gopher

type NotFound struct {
	Msg string
}

func (nf NotFound) Error() string {
	return nf.Msg
}
