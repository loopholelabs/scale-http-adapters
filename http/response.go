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
	"fmt"
	signature "github.com/loopholelabs/scale-signature-http"
	"net/http"
	"strings"
)

// FromResponse serializes the *ResponseWriter object into a runtime.Context
func FromResponse(ctx *signature.Context, w *ResponseWriter) {
	ctx.Response.StatusCode = int32(w.statusCode)

	if ctx.Response.Headers == nil {
		ctx.Response.Headers = signature.NewHttpResponseHeadersMap(uint32(len(w.headers)))
	}

	for k, v := range w.headers {
		ctx.Response.Headers[k] = &signature.HttpStringList{
			Value: v,
		}
	}
	ctx.Response.Body = w.buffer.Bytes()
}

// ToResponse deserializes the runtime.Context object into the http.ResponseWriter
func ToResponse(ctx *signature.Context, w http.ResponseWriter) error {
	for k, v := range ctx.Response.Headers {
		w.Header().Set(k, strings.Join(v.Value, ","))
	}

	if ctx.Response.StatusCode == 0 {
		w.WriteHeader(http.StatusOK)
	} else if ctx.Response.StatusCode < 100 || ctx.Response.StatusCode > 599 {
		return fmt.Errorf("invalid status code: %d", ctx.Response.StatusCode)
	} else {
		w.WriteHeader(int(ctx.Response.StatusCode))
	}

	if ctx.Response.Body != nil {
		_, err := w.Write(ctx.Response.Body)
		if err != nil {
			return fmt.Errorf("error writing response body: %w", err)
		}
	}

	return nil
}
