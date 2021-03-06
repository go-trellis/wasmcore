// GNU GPL v3 License
// Copyright (c) 2018 github.com:go-trellis

package servermsg

import "encoding/gob"

// RegisterTypes 注册
func RegisterTypes() {
	gob.Register(Queue{})
	gob.Register(Error{})
}

// Queue 队列
type Queue struct {
	Position int
	Done     bool
}

// Error 错误
type Error struct {
	Message string
}
