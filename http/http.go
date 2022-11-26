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

// Package http provides a Scale Runtime Adapter for the Standard net/http library.
package http

import (
	httpSignature "github.com/loopholelabs/scale-signature-http"
	runtime "github.com/loopholelabs/scale/go"
	"net/http"
)

// This compiler guard ensures that the HTTP adapter implements the net/http.Handler interface.
var _ http.Handler = (*HTTP)(nil)

// HTTP is a Scale Runtime Adapter for the Standard net/http library.
type HTTP struct {
	next    http.Handler
	runtime *runtime.Runtime[*httpSignature.Context]
}

// New returns a new HTTP adapter given a Scale Runtime and an optional net/http.Handler.
func New(next http.Handler, runtime *runtime.Runtime[*httpSignature.Context]) *HTTP {
	return &HTTP{
		next:    next,
		runtime: runtime,
	}
}

func (h *HTTP) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	i, err := h.runtime.Instance(h.Next(req))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = FromRequest(i.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = i.Run(req.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	err = ToResponse(i.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}

func (h *HTTP) Next(req *http.Request) runtime.Next[*httpSignature.Context] {
	if h.next == nil {
		return nil
	}
	return func(ctx *httpSignature.Context) (*httpSignature.Context, error) {
		ToRequest(ctx, req)
		w := NewResponseWriter()
		err := ToResponse(ctx, w)
		if err != nil {
			return nil, err
		}
		h.next.ServeHTTP(w, req)
		err = FromRequest(ctx, req)
		if err != nil {
			return nil, err
		}
		FromResponse(ctx, w)
		return ctx, nil
	}
}
