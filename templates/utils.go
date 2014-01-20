// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

/*  Filename:    utils.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-07-06 01:01:59.509179 -0700 PDT
 *  Description: 
 */

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
)

// Output file creation helper. Creates any missing parent directories. Does not
// overwrite existing files.
func FileCreate(path string) (*os.File, error) {
	if _, err := os.Stat(path); err == nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0644); err != nil { // TODO configurable? smarter?
		return nil, err
	}
	return os.Create(path)
}

// An artificial wrapper struct to make template rendering less cumbersome.
type Environment struct{ v interface{} }

// Create an Environment.
func Env(v interface{}) Environment { return Environment{v} }

// Render an ordered set of templates. Halts if a rendering error occurrs.
func (env Environment) Render(out io.Writer, ts Interface, names ...string) (err error) {
	for _, name := range names {
		if err = ts.Render(out, name, env.v); err != nil {
			return
		}
	}
	return
}

// Render a raw template string using the templates/functions of ts.
func (env Environment) RenderText(out io.Writer, ts Interface, prefix, text string) (err error) {
	randomness, err := randBase64(27)
	if err != nil {
		return
	}
	name := prefix + randomness
	if err = ts.Source(SourceTemplate{name, text}); err != nil {
		return
	}
	return env.Render(out, ts, name)
}

// Like Environment.RenderText, but returns a string containing rendered content.
func (env Environment) RenderTextAsString(ts Interface, prefix, text string) (string, error) {
	buff := new(bytes.Buffer)
	if err := env.RenderText(buff, ts, prefix, text); err != nil {
		return "", err
	}
	return buff.String(), nil
}

// A random string by URL-encoding n random bytes. Returned value length equal to (n*4+2)/3 (int division).
func randBase64(n int) (string, error) {
	source := make([]byte, n)
	if _, err := rand.Read(source); err != nil {
		return "", err
	}

	encLen := base64.URLEncoding.EncodedLen(n)
	buff := make([]byte, encLen)
	base64.URLEncoding.Encode(buff, source)
	return string(buff), nil
}
