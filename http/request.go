/*
	Copyright 2022 Loophole Labs

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

package http

import (
	"bytes"
	signature "github.com/loopholelabs/scale-signature-http"
	"io"
	"net/http"
)

const (
	BodyLimit = 1024 * 1024 * 10
)

// FromRequest serializes http.Request object into a runtime.Context
func FromRequest(ctx *signature.Context, req *http.Request) error {
	if ctx.Generated().Request.Headers == nil {
		ctx.Generated().Request.Headers = signature.NewHttpResponseHeadersMap(uint32(len(req.Header)))
	}
	for k, v := range req.Header {
		ctx.Generated().Request.Headers[k] = &signature.HttpStringList{
			Value: v,
		}
	}
	ctx.Generated().Request.Method = req.Method
	ctx.Generated().Request.ContentLength = req.ContentLength
	ctx.Generated().Request.Protocol = req.Proto
	ctx.Generated().Request.IP = req.RemoteAddr

	if req.ContentLength != 0 {
		var err error
		ctx.Generated().Request.Body, err = io.ReadAll(io.LimitReader(req.Body, BodyLimit))
		if err != nil {
			return err
		}
	} else {
		ctx.Generated().Request.Body = nil
	}

	return nil
}

// ToRequest deserializes the runtime.Context object into an existing http.Request
func ToRequest(ctx *signature.Context, req *http.Request) {
	req.Method = ctx.Generated().Request.Method
	req.ContentLength = ctx.Generated().Request.ContentLength
	req.Proto = ctx.Generated().Request.Protocol
	req.RemoteAddr = ctx.Generated().Request.IP

	for k, v := range ctx.Generated().Request.Headers {
		req.Header[k] = v.Value
	}

	if ctx.Generated().Request.ContentLength != 0 {
		req.Body = io.NopCloser(bytes.NewReader(ctx.Generated().Request.Body))
	} else {
		req.Body = io.NopCloser(bytes.NewReader(nil))
	}
}
