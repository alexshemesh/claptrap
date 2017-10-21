package contracts


type Market interface {
	GetBalance() (retVal Balance, err error)
}