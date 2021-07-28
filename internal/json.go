package internal

import (
	"github.com/hslatman/mud.yang.go/pkg/mudyang"
	"github.com/openconfig/ygot/ygot"
	"github.com/pkg/errors"
)

func JSON(mudfile *mudyang.Mudfile) (string, error) {
	json, err := ygot.EmitJSON(mudfile, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: true,
		},
		SkipValidation: false, // TODO: provide flag to skip?
	})
	if err != nil {
		return "", errors.Wrap(err, "could not marshal MUD file into JSON")
	}
	// TODO: ygot will alphabetically order properties (mapJSON); do we want some
	// way to override this, so that the informational parts can be shown in the top?
	// Of course this impacts ordering of the output, so care should be taken that
	// this behavior is known to the user.
	return json, nil
}
