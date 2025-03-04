package open_salt

import "golang.org/x/crypto/bcrypt"

func Encrypt(pwd string) (string, error) {
	if hashPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost); err != nil {
		return "", err
	} else {
		return string(hashPwd), nil
	}
}

func Valid(pwd string, real string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(real)); err != nil {
		return false
	}
	return true
}
