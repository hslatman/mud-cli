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
