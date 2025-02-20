package command

import "github.com/nalgeon/redka"

// Returns all fields and values in a hash.
// HGETALL key
// https://redis.io/commands/hgetall
type HGetAll struct {
	baseCmd
	key string
}

func parseHGetAll(b baseCmd) (*HGetAll, error) {
	cmd := &HGetAll{baseCmd: b}
	if len(cmd.args) != 1 {
		return cmd, ErrInvalidArgNum
	}
	cmd.key = string(cmd.args[0])
	return cmd, nil
}

func (cmd *HGetAll) Run(w Writer, red *redka.Tx) (any, error) {
	items, err := red.Hash().Items(cmd.key)
	if err != nil {
		w.WriteError(cmd.Error(err))
		return nil, err
	}
	w.WriteArray(len(items) * 2)
	for field, val := range items {
		w.WriteBulkString(field)
		w.WriteBulk(val)
	}
	return items, nil
}
