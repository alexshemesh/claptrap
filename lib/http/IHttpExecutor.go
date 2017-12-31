package httpClient

import "io"

type HttpExecutor interface {
	Get() (retVal HttpExecutor)
	Put() (retVal HttpExecutor)
	Post() (retVal HttpExecutor)
	Delete() (retVal HttpExecutor)
	SetContentType(newContentType string)( HttpClient)
	ExecuteWithReader(url string, headers map[string]string, at io.ReaderAt ) (response []byte, err error)
	Execute(url string, headers map[string]string, body []byte) (response []byte, err error)
}