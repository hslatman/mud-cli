package internal

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func Contents(filepath string) ([]byte, error) {
	if isURL(filepath) {
		// TODO: RFC states that MUD URL should be HTTPS; fail on that here if not?
		f, err := ioutil.TempFile(os.TempDir(), "")
		if err != nil {
			return []byte{}, err
		}
		defer f.Close()
		resp, err := http.Get(filepath)
		if err != nil {
			return []byte{}, err
		}
		defer resp.Body.Close()

		_, err = io.Copy(f, resp.Body)
		if err != nil {
			return []byte{}, err
		}
		filepath = f.Name()
	}

	json, err := ioutil.ReadFile(filepath)
	if err != nil {
		return []byte{}, fmt.Errorf("file could not be read: %w", err)
	}

	return json, nil
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	fmt.Println(err)
	fmt.Println(u)
	return err == nil && u.Scheme != "" && u.Host != ""
}
