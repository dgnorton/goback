package goback

import (
	"fmt"
	"testing"
)

func TestPush(t *testing.T) {
	txn := Transaction{}
	if len(txn.stack) != 0 {
		t.Errorf("Expected initial stack length of 0")
	}
	txn.Push(func() error {return fmt.Errorf("test")})
	if len(txn.stack) != 1 {
		t.Errorf("Expected stack length of 1 after Push")
	}
	txn.Push(func() error {return fmt.Errorf("test2")})
	if len(txn.stack) != 2 {
		t.Errorf("Expected stack length of 2 after second Push")
	}
}

func TestPop(t *testing.T) {
	txn := Transaction{}
	if len(txn.stack) != 0 {
		t.Errorf("Expected initial stack length of 0")
	}

	txn.Push(func() error {return fmt.Errorf("test")})
	if len(txn.stack) != 1 {
		t.Errorf("Expected stack length of 1 after Push")
	}

	txn.Push(func() error {return fmt.Errorf("test2")})
	if len(txn.stack) != 2 {
		t.Errorf("Expected stack length of 2 after second Push")
	}

	fn := txn.Pop()
	if fn == nil {
		t.Errorf("Expected function after first Pop but got nil")
	}
	if len(txn.stack) != 1 {
		t.Errorf("Expected stack length of 1 after first Pop")
	}
	err := fn()
	if err == nil {
		t.Errorf("Expected error from function")
	}
	etxt := err.Error()
	if etxt != "test2" {
		t.Errorf(`Expected error text to be "test2" but got "%s"`, etxt)
	}

	fn = txn.Pop()
	if fn == nil {
		t.Errorf("Expected function after second Pop")
	}
	if len(txn.stack) != 0 {
		t.Errorf("Expected stack length of 0 after second Pop")
	}
	err = fn()
	if err == nil {
		t.Error("Expected error from function")
	}
	etxt = err.Error()
	if etxt != "test" {
		t.Errorf(`Expected error text to be "test" but got "%s"`, etxt)
	}
}

func TestRollbackAfteCommit(t *testing.T) {
	txn := Transaction{}
	if len(txn.stack) != 0 {
		t.Errorf("Expected initial stack length of 0")
	}
	if txn.committed {
		t.Errorf("Expected initial state of txn.committed to be false")
	}

	txn.Push(func() error {return fmt.Errorf("test")})
	txn.Commit()

	if txn.committed == false {
		t.Errorf("Expected state of txn.committed == true after Commit")
	}

	err := txn.Rollback()

	if err != nil {
		t.Error(err)
	}
}

func TestRollback(t *testing.T) {
	txn := Transaction{}
	if len(txn.stack) != 0 {
		t.Errorf("Expected initial stack length of 0")
	}
	if txn.committed {
		t.Errorf("Expected initial state of txn.committed to be false")
	}

	cnt := 0

	txn.Push(func() error {cnt++; return nil})
	txn.Push(func() error {return fmt.Errorf("test2")})
	txn.Push(func() error {cnt++; return nil})
	txn.Push(func() error {cnt++; return nil})
	txn.Push(func() error {cnt++; return nil})

	if txn.committed {
		t.Errorf("Expected state of txn.committed to still be false")
	}

	err := txn.Rollback()

	if err == nil {
		t.Errorf("Expected err not to be nil after Rollback")
	}
	
	etxt := err.Error()
	if etxt != "test2" {
		t.Errorf(`Expected error text to be "test2" but got "%s"`, etxt)
	}

	if cnt != 3 {
		t.Errorf("Expected cnt == 3 but found cnt == %d", cnt)
	}
}