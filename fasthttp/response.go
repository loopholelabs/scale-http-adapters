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

func ToResponseContext(ctx *signature.Context, fastCTX *fasthttp.RequestCtx) {
	fastCTX.Response.SetStatusCode(int(ctx.Response.StatusCode))
	fastCTX.Response.SetBody(ctx.Response.Body)

	for k, v := range ctx.Response.Headers {
		fastCTX.Response.Header.Set(k, strings.Join(v.Value, ","))
	}
}

func FromResponseContext(ctx *signature.Context, fastCTX *fasthttp.RequestCtx) {
	ctx.Response.StatusCode = int32(fastCTX.Response.StatusCode())
	ctx.Response.Body = fastCTX.Response.Body()

	if ctx.Response.Headers == nil {
		ctx.Response.Headers = make(map[string]*signature.HttpStringList)
	}

	fastCTX.Response.Header.VisitAll(func(key []byte, value []byte) {
		ctx.Response.Headers[string(key)] = &signature.HttpStringList{
			Value: strings.Split(string(value), ","),
		}
	})
}
