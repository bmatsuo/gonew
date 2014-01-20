// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package extension

/*  Filename:    extensions.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-07-06 23:50:16.792611 -0700 PDT
 *  Description: 
 */

import ()

type Interface interface {
	Namespace() string
}

var Extensions map[string]interface{}

func Register(extension Interface) Interface {
	if Extensions == nil {
		Extensions = make(map[string]interface{})
	}
	Extensions[extension.Namespace()] = extension
	return extension
}
