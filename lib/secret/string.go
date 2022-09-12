package secret

type String string

func (s String) String() string {
	return "<secret>"
}
