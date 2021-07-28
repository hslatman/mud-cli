package internal

import (
	"github.com/hslatman/mud.yang.go/pkg/mudyang"
	"github.com/openconfig/ygot/ytypes"
)

func Validate(mudfile *mudyang.Mudfile) error {
	// TODO: more validation options?
	options := &ytypes.LeafrefOptions{
		IgnoreMissingData: false,
		Log:               true,
	}
	if err := mudfile.Validate(options); err != nil {
		return err
	}
	return nil
}
