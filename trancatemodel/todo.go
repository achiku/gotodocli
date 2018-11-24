package trancatemodel

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// ToDo todo table
type ToDo struct {
	ID          int64
	Description string
	Status      string
}

// Action action table
type Action struct {
	ID          int64
	ToDoID      int64
	Category    string
	PerformedAt time.Time
}

// Create creates todo
func (dt *ToDo) Create(q Queryer) error {
	err := q.QueryRow(
		`insert into todo (description, status) values ($1, $2) returning id`,
		dt.Description,
		dt.Status,
	).Scan(&dt.ID)
	if err != nil {
		return err
	}
	return nil
}

// Create creates action
func (ac *Action) Create(q Queryer) error {
	err := q.QueryRow(
		`insert into action (todo_id, category, performed_at) values ($1, $2, $3) returning id`,
		ac.ToDoID,
		ac.Category,
		ac.PerformedAt,
	).Scan(&ac.ID)
	if err != nil {
		return err
	}
	return nil
}

// CreateNewToDo creates new todo
func CreateNewToDo(tx Txer, td *ToDo) error {
	if err := td.Create(tx); err != nil {
		return errors.Wrap(err, "failed to create todo")
	}
	ac := Action{
		ToDoID:      td.ID,
		Category:    "created",
		PerformedAt: time.Now(),
	}
	if err := ac.Create(tx); err != nil {
		return errors.Wrap(err, "failed to create action")
	}
	return nil
}

// GetToDoByID get todo by id
func GetToDoByID(q Queryer, id int64) (*ToDo, bool, error) {
	var td ToDo
	err := q.QueryRow(`
	select
	  id
	  , description
	  , status
	from todo
	where id = $1
	`, id).Scan(
		&td.ID,
		&td.Description,
		&td.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &td, true, nil
}

// GetActionByToDoID get action by todo id
func GetActionByToDoID(q Queryer, id int64) ([]Action, error) {
	rows, err := q.Query(`
	select
	   id
	   , todo_id
	   , category
	   , performed_at
	from action
	where todo_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []Action
	for rows.Next() {
		var a Action
		rows.Scan(
			&a.ID,
			&a.Category,
			&a.PerformedAt,
		)
		as = append(as, a)
	}
	return as, nil
}
