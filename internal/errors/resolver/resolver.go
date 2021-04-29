// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resolver

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// errorResolvers is the list of known resolvers for kpt errors.
var errorResolvers []ErrorResolver

// AddErrorResolver adds the provided error resolver to the list of resolvers
// which will be used to resolve errors.
func AddErrorResolver(er ErrorResolver) {
	errorResolvers = append(errorResolvers, er)
}

// ResolveError attempts to resolve the provided error into a descriptive
// string which will be displayed to the user. If the last return value is false,
// the error could not be resolved.
func ResolveError(err error) (string, bool) {
	for _, resolver := range errorResolvers {
		msg, found := resolver.Resolve(err)
		if found {
			return msg, true
		}
	}
	return "", false
}

// ExecuteTemplate takes the provided template string and data, and renders
// the template. If something goes wrong, it panics.
func ExecuteTemplate(text string, data interface{}) (string, bool) {
	tmpl, tmplErr := template.New("kpterror").Parse(text)
	if tmplErr != nil {
		panic(fmt.Errorf("error creating template: %w", tmplErr))
	}

	var b bytes.Buffer
	execErr := tmpl.Execute(&b, data)
	if execErr != nil {
		panic(fmt.Errorf("error executing template: %w", execErr))
	}
	return strings.TrimSpace(b.String()), true
}

// ErrorResolver is an interface that allows kpt to resolve an error into
// an error message suitable for the end user.
type ErrorResolver interface {
	Resolve(err error) (string, bool)
}