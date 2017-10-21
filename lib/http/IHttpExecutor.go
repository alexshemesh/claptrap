package httpClient

type HttpExecutor interface {
	Get() (retVal HttpExecutor)
	Put() (retVal HttpExecutor)
	Post() (retVal HttpExecutor)
	Delete() (retVal HttpExecutor)
	SetContentType(newContentType string)( HttpClient)
	Execute(url string, headers map[string]string, body []byte) (response []byte, err error)
}