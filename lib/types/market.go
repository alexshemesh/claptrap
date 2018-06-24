package types


type Market interface {
	GetBalance() (retVal Balance, err error)
}