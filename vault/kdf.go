package vault

import (
	"fmt"
	"log"

	"github.com/alexedwards/argon2id"
	"golang.org/x/crypto/argon2"
)

func DeriveKey(password string) []byte {
	salt := []byte("fixed-or-random-salt-here")
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
}

func PasswordHashCheck() {
	password := "hunter2"
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hashed Password:", hash)

	passwordToCheck := []string{"hunter1", "hunter2", "hunter3"}

	for _, password := range passwordToCheck {
		match, err := argon2id.ComparePasswordAndHash(password, hash)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\"%s\":\t%vn", password, match)
	}

}
