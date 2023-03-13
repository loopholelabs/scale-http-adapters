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
	"strings"
)

const (
	BodyLimit = 1024 * 1024 * 10
)

// FromRequest serializes http.Request object into a runtime.Context
func FromRequest(ctx *signature.Context, req *http.Request) error {
	if ctx.Request.Headers == nil {
		ctx.Request.Headers = signature.NewHttpResponseHeadersMap(uint32(len(req.Header)))
	}
	for k, v := range req.Header {
		ctx.Request.Headers[strings.ToLower(k)] = &signature.HttpStringList{
			Value: v,
		}
	}
	ctx.Request.Method = req.Method
	ctx.Request.ContentLength = req.ContentLength
	ctx.Request.Protocol = req.Proto
	ctx.Request.IP = req.RemoteAddr

	if req.ContentLength != 0 {
		var err error
		ctx.Request.Body, err = io.ReadAll(io.LimitReader(req.Body, BodyLimit))
		if err != nil {
			return err
		}
	} else {
		ctx.Request.Body = nil
	}

	return nil
}

// ToRequest deserializes the runtime.Context object into an existing http.Request
func ToRequest(ctx *signature.Context, req *http.Request) {
	req.Method = ctx.Request.Method
	req.ContentLength = ctx.Request.ContentLength
	req.Proto = ctx.Request.Protocol
	req.RemoteAddr = ctx.Request.IP

	for k, v := range ctx.Request.Headers {
		req.Header[strings.ToLower(k)] = v.Value
	}

	if ctx.Request.ContentLength != 0 {
		req.Body = io.NopCloser(bytes.NewReader(ctx.Request.Body))
	} else {
		req.Body = io.NopCloser(bytes.NewReader(nil))
	}
}
