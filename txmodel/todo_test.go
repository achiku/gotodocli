package txmodel

import (
	"testing"
	"time"
)

func testShowCount(t *testing.T, tx Queryer) {
	t.Helper()
	var cnt int64
	if err := tx.QueryRow(`select count(*) from todo`).Scan(&cnt); err != nil {
		t.Fatal(err)
	}
	t.Logf("todo count=%d", cnt)

	if err := tx.QueryRow(`select count(*) from action`).Scan(&cnt); err != nil {
		t.Fatal(err)
	}
	t.Logf("action count=%d", cnt)
}

func TestGetToDoByID(t *testing.T) {
	tx, cleanup := TestSetupTx(t)
	defer cleanup()

	testShowCount(t, tx)
	td := TestCreateToDo(t, tx, &ToDo{
		Description: "test",
		Status:      "created",
	})

	cases := []struct {
		id    int64
		found bool
	}{
		{id: td.ID, found: true},
		{id: 0, found: false},
	}

	for _, c := range cases {
		td, found, err := GetToDoByID(tx, c.id)
		if err != nil {
			t.Fatal(err)
		}
		if found != c.found {
			t.Errorf("want %t got %t", c.found, found)
		}
		if found && td.ID != c.id {
			t.Errorf("want %d got %d", c.id, td.ID)
		}
	}
	testShowCount(t, tx)
}

func TestGetActionByToDoID(t *testing.T) {
	tx, cleanup := TestSetupTx(t)
	defer cleanup()

	testShowCount(t, tx)
	td := TestCreateToDo(t, tx, &ToDo{
		Description: "test",
		Status:      "created",
	})
	TestCreateAction(t, tx, &Action{
		ToDoID:      td.ID,
		Category:    "start",
		PerformedAt: time.Now(),
	})

	cases := []struct {
		id    int64
		count int
	}{
		{id: td.ID, count: 2},
		{id: 0, count: 0},
	}

	for _, c := range cases {
		acs, err := GetActionByToDoID(tx, c.id)
		if err != nil {
			t.Fatal(err)
		}
		if l := len(acs); l != c.count {
			t.Errorf("want %d got %d", c.count, l)
		}
	}
	testShowCount(t, tx)
}
