package goback

type TransactionFunc func() error

type Transaction struct {
	stack []TransactionFunc
	committed bool
}

func (t *Transaction) Push(f TransactionFunc) {
	t.stack = append(t.stack, f)
}

func (t *Transaction) Pop() TransactionFunc {
	ubnd := len(t.stack) - 1
	if ubnd < 0 {
		return nil
	}
	fn := t.stack[ubnd]
	t.stack = t.stack[:ubnd]
	return fn
}

func (t *Transaction) Commit() {
	t.committed = true
}

// Rollback all acctions by calling all rollback functions.
func (t *Transaction) Rollback() error {
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