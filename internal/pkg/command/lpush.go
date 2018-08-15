package command

import (
	"github.com/namreg/godown-v2/internal/pkg/storage"
)

func init() {
	cmd := new(Lpush)
	commands[cmd.Name()] = cmd
}

//Lpush is the LPUSH command
type Lpush struct{}

//Name implements Name of Command interface
func (c *Lpush) Name() string {
	return "LPUSH"
}

//Help implements Help of Command interface
func (c *Lpush) Help() string {
	return `Usage: LPUSH key value [value ...]
Prepend one or multiple values to a list.`
}

//Execute implements Execute of Command interface
func (c *Lpush) Execute(strg storage.Storage, args ...string) Result {
	if len(args) < 2 {
		return ErrResult{Value: ErrWrongArgsNumber}
	}

	strg.Lock()
	defer strg.Unlock()

	key := storage.Key(args[0])

	old, err := strg.Get(key)
	if err != nil && err != storage.ErrKeyNotExists {
		return ErrResult{Value: err}
	}

	vals := args[1:]

	// reverse vals
	for i, j := 0, len(vals)-1; i < j; i, j = i+1, j-1 {
		vals[i], vals[j] = vals[j], vals[i]
	}

	if old == nil {
		return c.put(strg, key, storage.NewListValue(vals))
	}

	if old.Type() != storage.ListDataType {
		return ErrResult{Value: ErrWrongTypeOp}
	}

	oldList := old.Data().([]string)

	newList := make([]string, 0, len(oldList)+len(vals))
	newList = append(newList, vals...)
	newList = append(newList, oldList...)

	return c.put(strg, key, storage.NewListValue(newList))
}

func (c *Lpush) put(strg storage.Storage, key storage.Key, value *storage.Value) Result {
	if err := strg.Put(key, value); err != nil {
		return ErrResult{Value: err}
	}
	return OkResult{}
}
