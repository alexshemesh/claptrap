package vault

import "fmt"

type VaultTestKit struct {
	values map[string]string
	errors map[string]error
}


func NewVaultTestKit()(retVal *VaultTestKit){
	retVal = &VaultTestKit{}
	retVal.values = make(map[string]string)
	retVal.errors = make(map[string]error)
	return retVal
}

func (this *VaultTestKit)SetValue(path string, value string) (err error){
	this.values[ path] = value
	return err
}

func (this VaultTestKit)GetValue(path string) (retVal string, err error){
	if this.values[ path] != "" {
		retVal = this.values[path]
		err = this.errors[path]
	}else if this.errors[path] != nil {
		err = this.errors[path]
	} else {
		err = fmt.Errorf("404 Not Found")
	}

	return retVal, err
}

func (this VaultTestKit)SetErrorEmulation(path string, value  error ) {
	this.errors[ path] = value
}
