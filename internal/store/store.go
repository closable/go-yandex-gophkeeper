package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/closable/go-yandex-gophkeeper/internal/errors"
	pb "github.com/closable/go-yandex-gophkeeper/internal/services/proto"
	_ "github.com/jackc/pgx/v5/stdlib"
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
	FileData struct {
		Data     string
		FilePath string
		DataType int
	}
)

// New new storage item
func New(connString string) (*Store, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	//PrepareDB(db)

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
	ctxIn, close := context.WithTimeout(context.Background(), timeout)
	defer close()

	tx, err := s.store.BeginTx(ctxIn, nil)
	if err != nil {
		return fmt.Errorf("%v %v", errors.ErrorTxDB, err)
	}
	defer tx.Rollback()

	sqlText := `INSERT INTO gophkeeper.users_data 
				(user_id, data_type, data, name) 
				VALUES($1, $2, $3, $4)`
	stmt, err := tx.PrepareContext(ctxIn, sqlText)
	if err != nil {
		return fmt.Errorf("%v %v", errors.ErrorPrepareDB, err)
	}

	_, err = stmt.ExecContext(ctxIn, userId, dataType, data, name)
	if err != nil {
		return fmt.Errorf("%v %v", errors.ErrorExecDB, err)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// UpdateItem update users data by id
func (s *Store) UpdateItem(userId, dataId int, data string) error {
	ctx, close := context.WithTimeout(context.Background(), time.Second*3)
	defer close()

	tx, err := s.store.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%v %v", errors.ErrorTxDB, err)
	}
	defer tx.Rollback()

	sqlText := "UPDATE gophkeeper.users_data SET data = $3 WHERE id = $2 and user_id = $1"
	stmt, err := tx.PrepareContext(ctx, sqlText)
	if err != nil {
		return fmt.Errorf("%v %v", errors.ErrorPrepareDB, err)
	}

	_, err = stmt.ExecContext(ctx, userId, dataId, data)
	if err != nil {
		return fmt.Errorf("%v %v", errors.ErrorExecDB, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%v %v", errors.ErrorCommitDB, err)
	}

	return nil
}

// DeleteItem delete users data by id
func (s *Store) DeleteItem(userId, dataId int) error {
	ctx, close := context.WithTimeout(context.Background(), time.Second*3)
	defer close()

	tx, err := s.store.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%v %v", errors.ErrorTxDB, err)
	}
	defer tx.Rollback()

	sqlText := `UPDATE gophkeeper.users_data SET is_deleted = true WHERE id = $2 and user_id = $1`
	stmt, err := tx.PrepareContext(ctx, sqlText)
	if err != nil {
		return fmt.Errorf("%v %v", errors.ErrorPrepareDB, err)
	}

	_, err = stmt.ExecContext(ctx, userId, dataId)
	if err != nil {
		return fmt.Errorf("%v %v", errors.ErrorExecDB, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%v %v", errors.ErrorCommitDB, err)
	}

	return nil
}

// ListItems list items by userId
func (s *Store) ListItems(userId int) ([]RowItem, error) {

	ctx, close := context.WithTimeout(context.Background(), time.Second*3)
	defer close()
	row := &RowItem{}
	items := make([]RowItem, 0)

	sqlText := `SELECT d.id, t.name, d.name, d.is_restore,
					case when d.data_type > 2 and substr(d.name, 1, 5) <> 'minio'  then 'данные файлов не отображаются' else d.data end data, length(d.data), d.data_type
					FROM gophkeeper.users_data d
					INNER JOIN gophkeeper.data_types t ON t.id = d.data_type
					WHERE d.user_id = $1 and not d.is_deleted
					ORDER BY d.id desc`
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

	if s.CheckUser(user) {
		return usr, fmt.Errorf("%v", errors.ErrorLoginInfo)
	}

	sqlText := `INSERT INTO gophkeeper.users 
				       (login, password, key) 
				VALUES ($1, sha256($2)::text, $3)`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	tx, err := s.store.BeginTx(ctx, nil)
	if err != nil {
		return usr, fmt.Errorf("%v %v", errors.ErrorTxDB, err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, sqlText)
	if err != nil {
		return usr, fmt.Errorf("%v %v", errors.ErrorPrepareDB, err)
	}

	_, err = stmt.ExecContext(ctx, user, pass, keyStr)
	if err != nil {
		return usr, fmt.Errorf("%v %v", errors.ErrorExecDB, err)
	}

	if err = tx.Commit(); err != nil {
		return usr, fmt.Errorf("%v %v", errors.ErrorCommitDB, err)
	}

	usr, err = s.GetUserInfo(user, pass)
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
func (s *Store) GetUserInfo(login, password string) (*UserDetail, error) {
	sqlText := `SELECT user_id, login, key 
				FROM gophkeeper.users 
				WHERE login = $1 and password = sha256($2)::text`
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	user := &UserDetail{}
	err := s.store.QueryRowContext(ctx, sqlText, login, password).Scan(&user.UserID, &user.Login, &user.KeyString)
	if err != nil {
		return user, fmt.Errorf("%v %v", errors.ErrorExecDB, err)
	}

	return user, nil
}

// GetUserKeyString helper functiion get key string by userID
func (s *Store) GetUserKeyString(userID int) (string, error) {
	sqlText := `SELECT key 
				FROM gophkeeper.users 
				WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var keyString string
	err := s.store.QueryRowContext(ctx, sqlText, userID).Scan(&keyString)
	if err != nil {
		return "", fmt.Errorf("%v %v", errors.ErrorExecDB, err)
	}

	return keyString, nil
}

// dropUser helper functions only for test
func (s *Store) DropUser(login string) error {
	sqlText := `DELETE FROM gophkeeper.users WHERE login = $1`
	_, err := s.store.Exec(sqlText, login)
	if err != nil {
		return fmt.Errorf("%v %v", errors.ErrorExecDB, err)
	}
	return nil
}

// dropData information helper functions only for test
func (s *Store) DropData(dataID int) error {
	sqlText := `DELETE FROM gophkeeper.users_data WHERE id = $1`
	_, err := s.store.Exec(sqlText, dataID)
	if err != nil {
		return err
	}
	return nil
}

// заглушка для удовлетворению интерфейса
func (s *Store) Upload(stream pb.FilseService_UploadServer) (*pb.FileUploadResponse, error) {
	var resp pb.FileUploadResponse
	return &resp, nil
}

// заглушка для удовлетворению интерфейса
func (s *Store) Download(in *pb.FileDownloadRequest, stream pb.FilseService_DownloadServer) error {
	return nil
}

// Download file data
func (s *Store) Health(n string) error {
	sqlText := "SELECT $1"
	var res int
	err := s.store.QueryRow(sqlText, n).Scan(&res)
	if err != nil {
		return err
	}
	return nil
}

// Download file data
func (s *Store) GetFileData(dataID int) (*FileData, error) {
	sqlText := `SELECT d.data, d.data_type, d.name file_path FROM gophkeeper.users_data d WHERE d.id = $1 and not d.is_deleted`

	data := &FileData{}
	err := s.store.QueryRow(sqlText, dataID).Scan(&data.Data, &data.DataType, &data.FilePath)
	if err != nil {
		return data, fmt.Errorf("%v %v", errors.ErrorExecDB, err)
	}

	return data, nil
}

// Сервисная функция, реализующая первоначальное состояние таблиц данных
// func PrepareDB(d *sql.DB) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
// 	defer cancel()
// 	pipe := make([]string, 5)
// 	pipe[0] = `CREATE SCHEMA IF NOT EXISTS gophkeeper`
// 	pipe[1] = `CREATE TABLE IF NOT EXISTS gophkeeper.users (
// 				user_id bigserial NOT NULL,
// 				login varchar(250) NOT NULL,
// 				password varchar(250) NOT NULL,
// 				is_active bool DEFAULT true NULL,
// 				key varchar(255) NULL,
// 				CONSTRAINT users_pkey PRIMARY KEY (user_id)
// 				)`
// 	pipe[2] = `CREATE TABLE IF NOT EXISTS gophkeeper.users_data (
// 				id bigserial NOT NULL,
// 				user_id bigserial NOT NULL,
// 				data_type int4 NULL,
// 				data text NULL,
// 				is_deleted bool DEFAULT false NULL,
// 				name varchar(255) NULL,
// 				is_restore bool DEFAULT false NULL,
// 				CONSTRAINT users_data_pkey PRIMARY KEY (id)
// 			)`
// 	pipe[3] = `CREATE TABLE IF NOT EXISTS gophkeeper.data_types (
// 				id bigserial NOT NULL,
// 				name varchar(255) NOT NULL,
// 				is_deleted bool DEFAULT false NULL,
// 				CONSTRAINT data_types_pkey PRIMARY KEY (id)
// 			)`
// 	pipe[4] = `merge into gophkeeper.data_types dt	using (
// 					select 1 id, 'PlainText' name,  false is_deleted
// 					union
// 					select 2 id, 'KeyValue' name, false is_deleted
// 					union
// 					select 3 id, 'FileData' name, false is_deleted
// 					union
// 					select 4 id, 'FolderData' name, false is_deleted
// 				) as res on (dt.id = res.id)
// 					when not matched then
// 					insert (id, name, is_deleted)
// 					values (res.id, res.name, res.is_deleted)`

// 	for ind, sql := range pipe {
// 		_, err := d.ExecContext(ctx, sql)
// 		if err != nil {
// 			fmt.Println(ind, sql, err)
// 			return err
// 		}
// 	}

// 	return nil
// }
