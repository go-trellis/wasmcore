// GNU GPL v3 License
// Copyright (c) 2018 github.com:go-trellis

package constormsg

import "encoding/gob"

func RegisterTypes() {
	gob.Register(Storing{})
}

type Storing struct {
	Starting  bool
	Finished  int
	Unchanged int
	Remain    int
	Done      bool
}
