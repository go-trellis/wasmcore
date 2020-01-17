// GNU GPL v3 License
// Copyright (c) 2018 github.com:go-trellis

package messages

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"

	"github.com/go-trellis/wasmcore"
	"github.com/go-trellis/wasmcore/server/constor/constormsg"
	"github.com/go-trellis/wasmcore/server/servermsg"
	"github.com/gorilla/websocket"
)

// Payload 上下文
type Payload struct {
	Message wasmcore.Message
}

func init() {
	// Commands:
	gob.Register(DeployQuery{})

	// Data messages:
	gob.Register(DeployQueryResponse{})
	gob.Register(DeployFileKey{})
	gob.Register(DeployFile{})
	gob.Register(DeployPayload{})
	gob.Register(DeployDone{})
	gob.Register(DeployClientVersionNotSupported{})

	// Initialise types in servermsg
	servermsg.RegisterTypes()

	// Initialise types in constormsg
	constormsg.RegisterTypes()
}

// Client sends a DeployQuery with all offered files.
// Server responds with DeployQueryResponse, with all required files listed.
// Client sends a DeployFile for each required file.

// DeployQuery 发布队列
type DeployQuery struct {
	Version string
	Files   []DeployFileKey
}

// DeployQueryResponse 返回队列
type DeployQueryResponse struct {
	Required []DeployFileKey
}

// DeployPayload 发布的上下文
type DeployPayload struct {
	Files []DeployFile
}

// DeployClientVersionNotSupported 不支持的客户端版本
type DeployClientVersionNotSupported struct{}

// DeployFileKey 发布文件的值
type DeployFileKey struct {
	Type DeployFileType
	Hash string // sha1 hash of contents
}

// DeployFile 发布文件的内容
type DeployFile struct {
	DeployFileKey
	Contents []byte // in the initial CommandDeploy, this is nil
}

// DeployDone 发布完成
type DeployDone struct{}

// DeployFileType 发布文件类型
type DeployFileType string

// 文件类型
const (
	DeployFileTypeIndex  DeployFileType = "index"
	DeployFileTypeLoader DeployFileType = "loader"
	DeployFileTypeWasm   DeployFileType = "wasm"
)

// Marshal 序列化
func Marshal(in wasmcore.Message) ([]byte, int, error) {
	p := Payload{in}
	buf := &bytes.Buffer{}
	gzw := gzip.NewWriter(buf)
	if err := gob.NewEncoder(gzw).Encode(p); err != nil {
		return nil, 0, err
	}
	if err := gzw.Close(); err != nil {
		return nil, 0, err
	}
	return buf.Bytes(), websocket.BinaryMessage, nil
}

// Unmarshal 反序列化
func Unmarshal(in []byte) (wasmcore.Message, error) {
	var p Payload
	gzr, err := gzip.NewReader(bytes.NewBuffer(in))
	if err != nil {
		return nil, err
	}
	if err := gob.NewDecoder(gzr).Decode(&p); err != nil {
		return nil, err
	}
	if err := gzr.Close(); err != nil {
		return nil, err
	}
	return p.Message, nil
}
