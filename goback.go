package goback

type TxFunc func() error

type Tx struct {
	stack []TxFunc
	committed bool
}

func Begin() *Tx {
	return &Tx{}
}

func (t *Tx) Exec(f TxFunc) {
	t.stack = append(t.stack, f)
}

func (t *Tx) Commit() {
	t.committed = true
}

// Rollback all acctions by calling all rollback functions.
func (t *Tx) Rollback() error {
	if t.committed {
		return nil
	}
	ubnd := len(t.stack) - 1
	for i, _ := range t.stack {
		err := t.stack[ubnd - i]()
		if err != nil {
			return err
		}
	}
	return nil
}