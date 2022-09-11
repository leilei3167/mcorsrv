package password

import "golang.org/x/crypto/bcrypt"

func Encode(s string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPwd), nil
}

func Compare(hashed string, inputPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(inputPwd))
	return err == nil
}
