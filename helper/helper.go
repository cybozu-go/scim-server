package helper

import "fmt"

type PhotoURL interface {
	Make(string, string) (string, error)
}

type NilPhotoURL struct{}

func (NilPhotoURL) Make(uid, path string) (string, error) {
	return fmt.Sprintf(`https://photo-storage/%s/%s`, uid, path), nil
}

type PhotoURLFunc func(string, string) (string, error)

func (f PhotoURLFunc) Make(uid, path string) (string, error) {
	return f(uid, path)
}
