# Modified the basic makefiles referred to from the
# Go home page.
#
# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG={{target}}
GOFILES=\
        {{main}}\

include $(GOROOT)/src/Make.pkg

