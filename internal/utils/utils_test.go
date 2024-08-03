package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestGetUserID(t *testing.T) {

	tests := []struct {
		name        string
		tokenString string
		userID      int
		want        bool
	}{
		// TODO: Add test cases.
		{
			name:        "Get UserID from auth token",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTI3MzgzOTAsIlVzZXJJRCI6NH0.V9WdWdJWeU1qqVCGDfTGu0asPZhiFUPmtnsfpN0GPro",
			userID:      4,
			want:        true,
		},
		{
			name:        "Get UserID from auth token ",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTI3MzgzOTAsIlVzZXJJRCI6NH0.V9WdWdJWeU1qqVCGDfTGu0asPZhiFUPmtnsfpN0GPro",
			userID:      5, // it can any int for test, excluding 4
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := GetUserID(tt.tokenString)
			if (userID == tt.userID) != tt.want {
				t.Errorf("GetUserID() = %v, compare with = %v,  want %v", userID, tt.userID, tt.want)
			}
		})
	}
}

func ExampleGetUserID() {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTI3MzgzOTAsIlVzZXJJRCI6NH0.V9WdWdJWeU1qqVCGDfTGu0asPZhiFUPmtnsfpN0GPro"
	//userID := 5
	out1 := GetUserID(tokenString)
	fmt.Println(out1)

	//Output
	//5

}

func TestCatStringData(t *testing.T) {
	s := "1234567890"
	wanted := "12345..."
	c := catStringData(s, 4)
	if c != wanted {
		t.Errorf("strings not equal %v %v", wanted, c)
	}

	c = catStringData(s, 20)
	if c != s {
		t.Errorf("strings must be equal %v %v", s, c)
	}
}

func TestMakePathFile(t *testing.T) {
	path := "cmd/gophkeeper/main.go"
	wanted := "cmd/gophkeeper/main.zip"
	s := MakePathFile(path, "zip")
	if s != wanted {
		t.Errorf("strings must be equal %v %v", wanted, path)
	}

	path = "t"
	wanted = ""
	s = MakePathFile(path, "zip")
	if s != wanted {
		t.Errorf("strings must be equal %v %v", wanted, path)
	}
}

func TestZipFolder(t *testing.T) {
	path := "../../cmd/gophkeeper/client/bin"
	wanted := "../../cmd/gophkeeper/client/bin.zip"
	output, err := ZipFolder(path)
	if err != nil {
		t.Errorf("ZipFolder %v", err)
	}

	if output != wanted {
		t.Errorf("Output error %v %v", wanted, output)
	}

	f, err := os.Open(wanted)

	if err != nil {
		t.Errorf("Zip file error %v", err)
	}

	stat, err := f.Stat()
	if err != nil {
		t.Errorf("File status error %v", err)
	}

	if stat.Size() < 100 {
		t.Errorf("Wrong file size %v", err)
	}

	err = os.Remove(wanted)
	if err != nil {
		t.Errorf("Wrong file name %v", err)
	}
}

func TestGetFileData(t *testing.T) {
	path := "../../cmd/gophkeeper/client/main.go"
	data, err := GetFileData(path)
	if err != nil {
		t.Errorf("Read file error %v", err)
	}

	if len(data) == 0 {
		t.Errorf("File must not be empty %v", err)
	}
}

func TestStoreFileData(t *testing.T) {
	path := "../../cmd/gophkeeper/client/test.txt"
	wanted := "testing process"
	err := StoreFileData(path, wanted)
	if err != nil {
		t.Errorf("File create error %v", err)
	}

	f, err := os.Open(path)

	if err != nil {
		t.Errorf("File error %v", err)
	}

	stat, err := f.Stat()
	if err != nil {
		t.Errorf("File status error %v", err)
	}

	if stat.Size() == 0 {
		t.Errorf("Wrong file size %v", err)
	}

	contents := make([]byte, stat.Size())
	_, err = f.Read(contents)
	if err != nil {
		t.Errorf("Wrong read file %v", err)
	}

	if string(contents) != wanted {
		t.Errorf("Contents must be equals %v %v", wanted, contents)
	}

	err = os.Remove(path)
	if err != nil {
		t.Errorf("Wrong file path %v", err)
	}
}

func TestEncrypt(t *testing.T) {
	keyString := GetRandomString(32)
	key, err := CryptoSeq(keyString)
	if err != nil {
		t.Errorf("CryptoSeq err %v", err)
	}

	wanted := "hello world"
	s := Encrypt(key, wanted)
	if len(s) == 0 {
		t.Errorf("Wrong encrypte result %v", s)
	}

	d := Decrypt(key, s)
	if d != wanted {
		t.Errorf("Result must be equal %v %v", wanted, d)
	}
}

func TestBuildJWTString(t *testing.T) {
	s, err := BuildJWTString(1)
	if err != nil {
		t.Errorf("Create JWT process error %v", err)
	}

	if len(s) == 0 {
		t.Errorf("Wrong JWT string %v", err)
	}
}
