package claymore

import (
	"testing"
	"path"
	"github.com/stretchr/testify/assert"

	"io/ioutil"
)

func Test_ParseStrainght(t *testing.T){
	fileContent,err := ioutil.ReadFile (path.Join("..","..","testdata","claymore_output.html"))

	assert.Nil(t,err)

	SplitTable(string(fileContent))
}


