package fw_test

import (
	"github.com/agorago/wego/fw"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestMakeCommandCatalog(t *testing.T) {
	cc := fw.MakeCommandCatalog()
	type someStruct struct{}
	ss := someStruct{}
	cc.RegisterCommand("ss",ss)
	assert.Equal(t,cc.Command("ss"),ss)
}