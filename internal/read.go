package internal

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/hslatman/mud.yang.go/pkg/mudyang"
	"github.com/pkg/errors"
)

func Read(filepath string) ([]byte, error) {
	if isURL(filepath) {
		// TODO: provide additional parameters for filetype to retrieve (or use a MUD client abstraction?)
		// so that the right headers can be provided in the requests (and check these in the response?)
		// TODO: RFC states that MUD URL should be HTTPS; fail on that here if not?
		// TODO: get rid of the temp file? I think it's relatively safe to keep the content bytes in memory
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
		if resp.StatusCode >= 400 {
			return []byte{}, fmt.Errorf("retrieving MUD file from %s resulted in HTTP status %d", filepath, resp.StatusCode)
		}
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

func Parse(data []byte) (*mudyang.Mudfile, error) {
	mud := &mudyang.Mudfile{}
	// TODO: provide options for unmarshaling?
	if err := mudyang.Unmarshal(data, mud); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal JSON")
	}
	return mud, nil
}

func ReadMUDFileFrom(filepath string) (*mudyang.Mudfile, error) {
	json, err := Read(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading file contents")
	}
	return Parse(json)
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
