goback
======

A simple non-database related transaction library

Typical use case
================

A function with a sequence of steps and if any of the steps fail all previous steps need to be reverted (e.g., variable values restored, file(s) deleted, etc.)

Example
=======

```
func Bar(i *int) {
	tx := goback.Begin()
	defer tx.Rollback()

	old := *i
	tx.Exec(func() error {*i = old; return nil})
	*i = 42

	err := fmt.Errorf("Forced error")
	
	if err == nil {
		tx.Commit()
	}
}

func main() {
	i := 10
	Bar(&i)
	fmt.Println(i)
}
```

Prints...

`10`

...because the `txn.Commit()` line never executes.
