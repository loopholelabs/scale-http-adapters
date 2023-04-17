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

package fasthttp

import (
	signature "github.com/loopholelabs/scale-signature-http"
	"github.com/valyala/fasthttp"
	"strings"
)

func FromRequestContext(ctx *signature.Context, fastCTX *fasthttp.RequestCtx) {
	ctx.Request.Protocol = "HTTP/1.1"
	ctx.Request.Method = string(fastCTX.Request.Header.Method())
	ctx.Request.IP = fastCTX.RemoteAddr().String()
	ctx.Request.ContentLength = int64(fastCTX.Request.Header.ContentLength())
	ctx.Request.Body = fastCTX.Request.Body()
	ctx.Request.URI = fastCTX.Request.URI().String()
	if ctx.Request.ContentLength < 1 {
		ctx.Request.ContentLength = int64(len(ctx.Request.Body))
	}

	if ctx.Request.Headers == nil {
		ctx.Request.Headers = make(map[string]*signature.HttpStringList)
	}

	fastCTX.Request.Header.VisitAll(func(key []byte, value []byte) {
		ctx.Request.Headers[strings.ToLower(string(key))] = &signature.HttpStringList{
			Value: strings.Split(string(value), ","),
		}
	})
}

func ToRequestContext(ctx *signature.Context, fastCTX *fasthttp.RequestCtx) {
	fastCTX.Request.Header.SetMethod(ctx.Request.Method)
	fastCTX.Request.Header.SetContentLength(int(ctx.Request.ContentLength))
	fastCTX.Request.SetBody(ctx.Request.Body)

	for k, v := range ctx.Request.Headers {
		fastCTX.Request.Header.Set(strings.ToLower(k), strings.Join(v.Value, ","))
	}
}
