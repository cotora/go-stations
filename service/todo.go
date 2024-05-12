package service

import (
	"context"
	"database/sql"
	"time"
	"fmt"
	"strings"
	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	res_insert,err:=s.db.ExecContext(ctx,insert,subject,description)
	if err!=nil{
		return nil,err
	}

	id,err:=res_insert.LastInsertId()
	if err!=nil{
		return nil,err
	}

	row:=s.db.QueryRowContext(ctx,confirm,id)

	var res_subject,res_description string
	var created_at_str,updated_at_str string
	err=row.Scan(&res_subject,&res_description,&created_at_str,&updated_at_str)
	if err!=nil{
		return nil,err
	}

	var created_at,updated_at time.Time
	created_at,_= time.Parse("2006-01-02T15:04:05Z07:00", created_at_str)
	updated_at,_= time.Parse("2006-01-02T15:04:05Z07:00", updated_at_str)

	res:= &model.TODO{
		ID: id,
		Subject: res_subject,
		Description: res_description,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	}
	return res, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	var rows *sql.Rows

	if size == 0{
		size=5
	}

	if prevID == 0{
		var err error
		rows, err = s.db.QueryContext(ctx, read, size)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
		if err != nil {
			return nil, err
		}
	}

	var todos []*model.TODO = make([]*model.TODO, 0)

	for rows.Next() {
		var id int64
		var subject, description, created_at_str, updated_at_str string
		err := rows.Scan(&id, &subject, &description, &created_at_str, &updated_at_str)
		if err != nil {
			return nil, err
		}
		created_at,_:= time.Parse("2006-01-02T15:04:05Z07:00", created_at_str)
		updated_at,_:= time.Parse("2006-01-02T15:04:05Z07:00", updated_at_str)

		todo := &model.TODO{
			ID:id,
			Subject: subject,
			Description: description,
			CreatedAt: created_at,
			UpdatedAt: updated_at,
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	res_update,err:=s.db.ExecContext(ctx,update,subject,description,id)
	if err!=nil{
		return nil,err
	}

	affected_num,err:=res_update.RowsAffected()
	if err!=nil{
		return nil,err
	}

	if affected_num==0{
		return nil,&model.ErrNotFound{}
	}

	id_update,err:=res_update.LastInsertId()
	if err!=nil{
		return nil,err
	}

	row:=s.db.QueryRowContext(ctx,confirm,id_update)

	var res_subject,res_description string
	var created_at_str,updated_at_str string
	err=row.Scan(&res_subject,&res_description,&created_at_str,&updated_at_str)
	if err!=nil{
		return nil,err
	}

	var created_at,updated_at time.Time
	created_at,_= time.Parse("2006-01-02T15:04:05Z07:00", created_at_str)
	updated_at,_= time.Parse("2006-01-02T15:04:05Z07:00", updated_at_str)

	res:= &model.TODO{
		ID: id_update,
		Subject: res_subject,
		Description: res_description,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	}

	return res, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	if len(ids)==0{
		return nil
	}

	query:=fmt.Sprintf(deleteFmt,strings.Repeat(",?",len(ids)-1))

	var tmp []interface{}=make([]interface{},len(ids))

	for i,id:=range ids{
		tmp[i]=id
	}

	res_delete,err:=s.db.ExecContext(ctx,query,tmp...)
	if err!=nil{
		return err
	}

	affected_num,err:=res_delete.RowsAffected()
	if err!=nil{
		return err
	}
	
	if affected_num==0{
		return &model.ErrNotFound{}
	}

	return nil
}
