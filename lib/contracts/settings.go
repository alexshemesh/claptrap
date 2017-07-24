package contracts

type Settings interface {
	GetValue(path string)(retVal string, err error)
}