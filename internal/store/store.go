package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/closable/go-yandex-gophkeeper/internal/utils"
)

type (
	Store struct {
		store *sql.DB
	}

	RowItem struct {
		Id        int
		Type      string
		Name      string
		IsRestore bool
		EncData   string
		Length    int
		DataType  int
	}
	UserDetail struct {
		UserID    int
		Login     string
		KeyString string
	}
)

// New new storage item
func New(connString string) (*Store, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	return &Store{
		store: db,
	}, nil
}

// AddItem add item into store
func (s *Store) AddItem(userId int, dataType int, data, name string) error {

	timeout := time.Millisecond * 100
	if dataType > 2 {
		timeout = time.Second * 10
	}
	ctx, close := context.WithTimeout(context.Background(), timeout)
	defer close()

	tx, err := s.store.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sqlText := `INSERT INTO gophkeeper.users_data 
				(user_id, data_type, data, name) 
				VALUES($1, $2, $3, $4)`
	stmt, err := tx.PrepareContext(ctx, sqlText)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, userId, dataType, data, name)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// UpdateItem update users data by id
func (s *Store) UpdateItem(userId, dataId int, data string) error {
	ctx, close := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer close()

	tx, err := s.store.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sqlText := "UPDATE gophkeeper.users_data SET data = $3 WHERE id = $2 and user_id = $1"
	stmt, err := tx.PrepareContext(ctx, sqlText)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, userId, dataId, data)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// DeleteItem delete users data by id
func (s *Store) DeleteItem(userId, dataId int) error {
	ctx, close := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer close()

	tx, err := s.store.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sqlText := `UPDATE gophkeeper.users_data SET is_deleted = true WHERE id = $2 and user_id = $1`
	stmt, err := tx.PrepareContext(ctx, sqlText)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, userId, dataId)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// ListItems list items by userId
func (s *Store) ListItems(userId int) ([]RowItem, error) {

	ctx, close := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer close()
	row := &RowItem{}
	items := make([]RowItem, 0)

	sqlText := `SELECT d.id, t.name, d.name, d.is_restore,
					case when d.data_type > 2 then 'данные файлов не отображаются' else d.data end data, length(d.data), d.data_type
					FROM gophkeeper.users_data d
					INNER JOIN gophkeeper.data_types t ON t.id = d.data_type
					WHERE d.user_id = $1 and not d.is_deleted`
	rows, err := s.store.QueryContext(ctx, sqlText, userId)
	if err != nil || rows.Err() != nil {
		return items, err
	}

	for rows.Next() {
		err = rows.Scan(&row.Id, &row.Type, &row.Name, &row.IsRestore, &row.EncData, &row.Length, &row.DataType)
		if err != nil {
			return items, err
		}
		items = append(items, *row)
	}

	return items, nil
}

// CreateUser create new user
func (s *Store) CreateUser(user, pass, keyStr string) (*UserDetail, error) {
	usr := &UserDetail{}
	userDecoded := utils.Decrypt(keyStr, user)
	passDecoded := utils.Decrypt(keyStr, pass)

	if s.CheckUser(userDecoded) {
		return usr, fmt.Errorf("login alread busy ")
	}

	sqlText := `INSERT INTO gophkeeper.users 
				       (login, password, key) 
				VALUES ($1, sha256($2)::text, $3)`
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	tx, err := s.store.BeginTx(ctx, nil)
	if err != nil {
		return usr, err
		//return fmt.Errorf("%s %w", errors_api.ErrorBeginTx.Error(), err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, sqlText)
	if err != nil {
		return usr, err
	}

	_, err = stmt.ExecContext(ctx, userDecoded, passDecoded, keyStr)
	if err != nil {
		return usr, err
	}

	if err = tx.Commit(); err != nil {
		return usr, err
	}

	usr, err = s.GetUserInfo(userDecoded)
	if err != nil {
		return usr, err
	}

	return usr, nil
}

// CheckUser check user name
func (s *Store) CheckUser(user string) bool {
	sqlText := "SELECT count(*) FROM gophkeeper.users WHERE login = $1"
	cnt := 0
	err := s.store.QueryRow(sqlText, user).Scan(&cnt)
	if err != nil {
		return true
	}
	return cnt > 0
}

// GetUserInfo get user detail info by login
func (s *Store) GetUserInfo(login string) (*UserDetail, error) {
	sqlText := `SELECT user_id, login, key FROM gophkeeper.users WHERE login = $1`
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	user := &UserDetail{}
	err := s.store.QueryRowContext(ctx, sqlText, login).Scan(&user.UserID, &user.Login, &user.KeyString)
	if err != nil {
		return user, err
	}

	return user, nil
}

/*
CREATE SCHEMA gophkeeper
    AUTHORIZATION postgres;
-----------------------------------
CREATE TABLE gophkeeper.users
(
    user_id bigserial NOT NULL,
    login character varying(250) NOT NULL,
    password character varying(250) NOT NULL,
    is_active boolean DEFAULT true,
    PRIMARY KEY (user_id)
);

ALTER TABLE IF EXISTS gophkeeper.users
    OWNER to postgres;
-----------------------------------
CREATE TABLE gophkeeper.users_data
(
    id bigserial NOT NULL,
    user_id bigserial NOT NULL,
    data_type integer,
    data text,
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS gophkeeper.users_data
    OWNER to postgres;
-----------------------------------

CREATE TABLE gophkeeper.data_types
(
    id bigserial NOT NULL,
    name character varying(255) NOT NULL,
    is_deleted boolean DEFAULT false,
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS gophkeeper.data_types
    OWNER to postgres;
-----------------------------------

*/
