package internal

import (
	"net/url"
	fp "path/filepath"

	"github.com/hslatman/mud.yang.go/pkg/mudyang"
)

func SignatureFilename(filepath string) (string, error) {
	var filename string
	if isURL(filepath) {
		u, err := url.Parse(filepath)
		if err != nil {
			return "", err
		}
		filename = fp.Base(u.EscapedPath())
		return filename + ".p7s", nil
		//u.Path = filename + ".p7s"
		//return u.String(), nil
	}

	// TODO: we probably need to change this for local files.
	// We actually need to know the location the signature will
	// be put online, so that that can be used as the signature
	// path in the MUD file and then can be signed.
	filename = fp.Base(filepath) + ".p7s"
	//return fp.Join(fp.Dir(filepath), filename), nil
	return filename, nil
}

func MUDURL(mudfile *mudyang.Mudfile) (*url.URL, error) {
	return url.Parse(*mudfile.Mud.MudUrl)
}

func MUDSignatureURL(mudfile *mudyang.Mudfile) (*url.URL, error) {
	sig := mudfile.Mud.MudSignature
	if sig == nil {
		return nil, nil
	}
	return url.Parse(*sig)
}

func NewMUDSignatureURL(mudurl *url.URL, filename string) *url.URL {
	sigURL := *mudurl
	path := sigURL.EscapedPath()
	mudFileName := fp.Base(path)
	path = path[0:len(path)-len(mudFileName)] + filename
	sigURL.Path = path
	return &sigURL
}

// func SignaturePath(filepath string) (string, error) {
// 	var filename string
// 	if isURL(filepath) {
// 		u, err := url.Parse(filepath)
// 		fmt.Println(u)
// 		if err != nil {
// 			return "", err
// 		}
// 		filename = u.EscapedPath()
// 		filename = fp.Base(filename)
// 		fmt.Println(filename)
// 	} else {
// 		filename = fp.Base(filepath)
// 	}

// 	extension := fp.Ext(filename)
// 	fmt.Println(extension)
// 	//filename = fp.Ext(filepath)
// 	name := filename[0 : len(filename)-len(extension)]

// 	return name + ".p7s", nil
// }
