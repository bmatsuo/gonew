// Copyright 2014, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// strings_test.go [created: Sat, 18 Jan 2014]

package extension


import "testing"

func TestStrings(t *testing.T) {
	s := String{}
	uc, err := s.UpperCamel("testing-testing")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if uc != "TestingTesting" {
		t.Errorf("unexpected UpperCamel: %v", uc)
	}
}
