package command

import (
	"testing"

	"github.com/nalgeon/redka/internal/testx"
)

func TestIncrByParse(t *testing.T) {
	tests := []struct {
		name string
		args [][]byte
		want IncrBy
		err  error
	}{
		{
			name: "incrby",
			args: buildArgs("incrby"),
			want: IncrBy{},
			err:  ErrInvalidArgNum,
		},
		{
			name: "incrby age",
			args: buildArgs("incrby", "age"),
			want: IncrBy{},
			err:  ErrInvalidArgNum,
		},
		{
			name: "incrby age 42",
			args: buildArgs("incrby", "age", "42"),
			want: IncrBy{key: "age", delta: 42},
			err:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd, err := Parse(test.args)
			testx.AssertEqual(t, err, test.err)
			if err == nil {
				cm := cmd.(*IncrBy)
				testx.AssertEqual(t, cm.key, test.want.key)
				testx.AssertEqual(t, cm.delta, test.want.delta)
			}
		})
	}
}

func TestIncrByExec(t *testing.T) {
	db, tx := getDB(t)
	defer db.Close()

	t.Run("create", func(t *testing.T) {
		cmd := mustParse[*IncrBy]("incrby age 42")
		conn := new(fakeConn)
		res, err := cmd.Run(conn, tx)
		testx.AssertNoErr(t, err)
		testx.AssertEqual(t, res, 42)
		testx.AssertEqual(t, conn.out(), "42")

		age, _ := db.Str().Get("age")
		testx.AssertEqual(t, age.MustInt(), 42)
	})

	t.Run("incrby", func(t *testing.T) {
		_ = db.Str().Set("age", "25")

		cmd := mustParse[*IncrBy]("incrby age 42")
		conn := new(fakeConn)
		res, err := cmd.Run(conn, tx)
		testx.AssertNoErr(t, err)
		testx.AssertEqual(t, res, 67)
		testx.AssertEqual(t, conn.out(), "67")

		age, _ := db.Str().Get("age")
		testx.AssertEqual(t, age.MustInt(), 67)
	})
}
