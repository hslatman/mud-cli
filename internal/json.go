/*
Copyright © 2021 Herman Slatman <hslatman>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
