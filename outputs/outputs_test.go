package outputs

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func Test_Outputers(t *testing.T) {
	outputer := GetOutputer("json")
	assert.Equal(t, "json", outputer.Name())

	registeredOutputers := []string{
		"documentation",
		"json",
		"json_oneline",
		"junit",
		"nagios",
		"rspecish",
		"silent",
		"tap",
	}
	assert.Equal(t, registeredOutputers, Outputers())
}
