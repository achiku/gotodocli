package txmodel

import (
	"testing"
	"time"
)

// TestCreateToDo creates test todo data
func TestCreateToDo(t *testing.T, q Queryer, td *ToDo) *ToDo {
	if err := td.Create(q); err != nil {
		t.Fatal(err)
	}
	ac := Action{
		ToDoID:      td.ID,
		Category:    "created",
		PerformedAt: time.Now(),
	}
	if err := ac.Create(q); err != nil {
		t.Fatal(err)
	}
	return td
}

// TestCreateAction creates test action data
func TestCreateAction(t *testing.T, q Queryer, ac *Action) *Action {
	if err := ac.Create(q); err != nil {
		t.Fatal(err)
	}
	return ac
}
