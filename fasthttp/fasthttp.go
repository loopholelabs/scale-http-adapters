/*
	Copyright 2023 Loophole Labs

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

// Package fasthttp provides a Scale Runtime Adapter for the fasthttp library
package fasthttp

import (
	signature "github.com/loopholelabs/scale-signature-http"
	scale "github.com/loopholelabs/scale/go"
	"github.com/valyala/fasthttp"
)

type FastHTTP struct {
	runtime *scale.Runtime[*signature.Context]
	next    fasthttp.RequestHandler
}

func New(next fasthttp.RequestHandler, runtime *scale.Runtime[*signature.Context]) *FastHTTP {
	return &FastHTTP{
		runtime: runtime,
		next:    next,
	}
}

func (f *FastHTTP) Handle(ctx *fasthttp.RequestCtx) {
	i, err := f.runtime.Instance(f.Next(ctx))
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	FromRequestContext(i.Context(), ctx)
	err = i.Run(ctx)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadGateway)
		return
	}
	ToResponseContext(i.Context(), ctx)
}

func (f *FastHTTP) Next(fastCTX *fasthttp.RequestCtx) scale.Next[*signature.Context] {
	if f.next == nil {
		return nil
	}
	return func(ctx *signature.Context) (*signature.Context, error) {
		ToRequestContext(ctx, fastCTX)
		ToResponseContext(ctx, fastCTX)
		f.next(fastCTX)
		FromRequestContext(ctx, fastCTX)
		FromResponseContext(ctx, fastCTX)
		return ctx, nil
	}
}
