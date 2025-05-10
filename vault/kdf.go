package vault

import (
	"fmt"
	"log"

	"github.com/alexedwards/argon2id"
)

var Parameters = &argon2id.Params{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

// Hash Master Password takes users master password to enter the app and returns an encoded argon2id string
func HashMasterPassword(password string) (encodedHash string, err error) {
	return argon2id.CreateHash(password, Parameters)
}

func VerifyMasterPassword(password, encodedHash string) (match bool, err error) {
	return argon2id.ComparePasswordAndHash(password, encodedHash)
}

// this is a test func to see how the package works
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
