package command

import (
	"github.com/namreg/godown-v2/internal/pkg/storage"
)

func init() {
	cmd := new(Hget)
	commands[cmd.Name()] = cmd
}

//Hget is the HGET command
type Hget struct{}

//Name implements Name of Command interface
func (c *Hget) Name() string {
	return "HGET"
}

//Help implements Help of Command interface
func (c *Hget) Help() string {
	return `Usage: HGET key field
Returns the value associated with field in the hash stored at key.`
}

//Execute implements Execute of Command interface
func (c *Hget) Execute(strg storage.Storage, args ...string) Result {
	if len(args) != 2 {
		return ErrResult{Value: ErrWrongArgsNumber}
	}

	strg.RLock()
	value, err := strg.Get(storage.Key(args[0]))
	strg.RUnlock()
	if err != nil {
		if err == storage.ErrKeyNotExists {
			return NilResult{}
		}
		return ErrResult{Value: err}
	}
	if value.Type() != storage.MapDataType {
		return ErrResult{Value: ErrWrongTypeOp}
	}
	m := value.Data().(map[string]string)
	if v, ok := m[args[1]]; ok {
		return StringResult{Value: v}
	}
	return NilResult{}
}
