package exception

type NotFound struct {
	Msg string
}

func (nf NotFound) Error() string {
	return nf.Msg
}
