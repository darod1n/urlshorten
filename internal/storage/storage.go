package storage

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

func GenerateToken(lenToken int) string {
	var token string
	rand.Seed(time.Now().UnixNano())
	alphabet := GetStrFromFile("alphabet")
	for i := 0; i < lenToken; i++ {
		r := rand.Intn(len(alphabet))
		token += string(alphabet[r])
	}
	return token
}

func GetStrFromFile(path string) string {
	var str string
	dir, _ := os.Getwd()
	file, err := os.Open(dir + "/secrets/" + path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	data := make([]byte, 64)
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		str += string(data[:n])
	}

	return str
}

func SaveToFile(path string, str string) {
	dir, _ := os.Getwd()
	file, err := os.Create(dir + "/secrets/" + path)
	if err != nil {
		os.Exit(1)
	}

	defer file.Close()
	file.WriteString(str)
}

func IsExist(token string) bool {
	dir, _ := os.Getwd()
	if _, err := os.Stat(dir + "/secrets/" + token); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
