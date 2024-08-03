package store

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/closable/go-yandex-gophkeeper/internal/config"
)

func TestCreateUser(t *testing.T) {
	cfg := config.LoadConfig()

	if len(cfg.DSN) == 0 {
		t.Errorf("Wrong config %v", cfg.DSN)
	}

	st, err := New(cfg.DSN)
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name     string
		login    string
		password string
	}{
		{
			name:     "create user 1",
			login:    fmt.Sprintf("user_%d", rand.Intn(100)),
			password: "test1",
		},
		{
			name:     "create user 2",
			login:    fmt.Sprintf("user_%d", rand.Intn(100)),
			password: "test2",
		},
		{
			name:     "create user 3",
			login:    fmt.Sprintf("user_%d", rand.Intn(100)),
			password: "test3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail, err := st.CreateUser(tt.login, tt.password, "")
			if err != nil {
				t.Errorf("Create user error %v %v", tt.name, err)
			}

			if detail.Login != tt.login {
				t.Errorf("Result must be equal %v %v %v", tt.name, detail.Login, tt.login)
			}

			s, err := st.GetUserKeyString(detail.UserID)
			if err != nil {
				t.Errorf("Wrong KeyString %v %v", tt.name, s)
			}

			err = st.AddItem(detail.UserID, 1, "123-456-789 00", "bank card")
			if err != nil {
				t.Errorf("Error add information %v %v", tt.name, err)
			}

			list, err := st.ListItems(detail.UserID)
			if err != nil {
				t.Errorf("Error list information %v %v", tt.name, err)
			}

			err = st.UpdateItem(detail.UserID, list[0].Id, "999-999-999 99")
			if err != nil {
				t.Errorf("Error update information %v %v", tt.name, err)
			}

			err = st.DeleteItem(detail.UserID, list[0].Id)
			if err != nil {
				t.Errorf("Error delete from list %v %v", tt.name, err)
			}

			err = st.DropData(list[0].Id)
			if err != nil {
				t.Errorf("Error data delete %v %v", tt.name, err)
			}

			err = st.DropUser(tt.login)
			if err != nil {
				t.Errorf("Invalid operations (delete user in db) %v %v %v", tt.name, detail.Login, tt.login)
			}
		})
	}

}
