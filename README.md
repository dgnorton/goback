goback
======

A simple non-database related transaction library

Example
=======

```
func Bar(i *int) {
	var txn goback.Transaction
	defer txn.Rollback()

	old := *i
	txn.Push(func() error {*i = old; return nil})
	*i = 42
	
	//txn.Commit()
}

func main() {
	i := 10
	Bar(&i)
	fmt.Println(i)
}
```
