package goback

import (
	"fmt"
	"testing"
)

func TestExec(t *testing.T) {
	tx := Begin()
	if len(tx.stack) != 0 {
		t.Errorf("Expected initial stack length of 0")
	}
	tx.Exec(func() error {return fmt.Errorf("test")})
	if len(tx.stack) != 1 {
		t.Errorf("Expected stack length of 1 after Exec")
	}
	tx.Exec(func() error {return fmt.Errorf("test2")})
	if len(tx.stack) != 2 {
		t.Errorf("Expected stack length of 2 after second Exec")
	}
}

func TestRollbackAfteCommit(t *testing.T) {
	tx := Begin()
	if len(tx.stack) != 0 {
		t.Errorf("Expected initial stack length of 0")
	}
	if tx.committed {
		t.Errorf("Expected initial state of tx.committed to be false")
	}

	tx.Exec(func() error {return fmt.Errorf("test")})
	tx.Commit()

	if tx.committed == false {
		t.Errorf("Expected state of tx.committed == true after Commit")
	}

	err := tx.Rollback()

	if err != nil {
		t.Error(err)
	}
}

func TestRollback(t *testing.T) {
	tx := Begin()
	if len(tx.stack) != 0 {
		t.Errorf("Expected initial stack length of 0")
	}
	if tx.committed {
		t.Errorf("Expected initial state of tx.committed to be false")
	}

	cnt := 0

	tx.Exec(func() error {cnt++; return nil})
	tx.Exec(func() error {return fmt.Errorf("test2")})
	tx.Exec(func() error {cnt++; return nil})
	tx.Exec(func() error {cnt++; return nil})
	tx.Exec(func() error {cnt++; return nil})

	if tx.committed {
		t.Errorf("Expected state of tx.committed to still be false")
	}

	err := tx.Rollback()

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