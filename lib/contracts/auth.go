package contracts



type Auth interface {
	LogIn( userName string, password string ) ( retObjSettings Settings,retObjAuth Auth ,err error)
	LogInCached(userName string)(retObjSettings Settings,retObjAuth Auth ,err error)

	SaveCachedToken(userName string , token string) (err error)

	GetUserName()(retVal string)
	GetUserToken()(retVal string)

	RenewToken()(err error)
}

