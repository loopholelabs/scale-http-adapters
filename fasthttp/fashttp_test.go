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
	"context"
	scale "github.com/loopholelabs/scale/go"
	"github.com/loopholelabs/scale/go/tests/harness"
	"github.com/loopholelabs/scalefile/scalefunc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

type TestCase struct {
	Name   string
	Module *harness.Module
	Run    func(*scalefunc.ScaleFunc, *testing.T)
}

func TestFastHTTP(t *testing.T) {
	nextModule := &harness.Module{
		Name:      "next",
		Path:      "../tests/modules/next/next.go",
		Signature: "github.com/loopholelabs/scale-signature-http",
	}

	modules := []*harness.Module{nextModule}

	generatedModules := harness.GoSetup(t, modules, "github.com/loopholelabs/scale-http-adapters/tests/modules")

	var testCases = []TestCase{
		{
			Name:   "Next",
			Module: nextModule,
			Run: func(scaleFunc *scalefunc.ScaleFunc, t *testing.T) {
				r, err := scale.New(context.Background(), []*scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				adapter := New(func(ctx *fasthttp.RequestCtx) {
					require.Equal(t, "GET", string(ctx.Method()))
					assert.Equal(t, "Test Data", string(ctx.PostBody()))
					ctx.Response.SetBody(append(ctx.Response.Body(), "-next"...))
					ctx.Response.SetStatusCode(200)
				}, r)

				s := &fasthttp.Server{
					Handler:         adapter.Handle,
					CloseOnShutdown: true,
					IdleTimeout:     time.Second,
				}

				ln := fasthttputil.NewInmemoryListener()
				var wg sync.WaitGroup

				wg.Add(1)
				go func() {
					defer wg.Done()
					require.NoError(t, s.Serve(ln))
				}()

				req, err := http.NewRequest("GET", "http://localhost:8080", strings.NewReader("Test Data"))
				assert.NoError(t, err)

				client := &http.Client{
					Transport: &http.Transport{
						DialContext: func(_ context.Context, _ string, _ string) (net.Conn, error) {
							return ln.Dial()
						},
					},
				}

				res, err := client.Do(req)
				assert.NoError(t, err)

				body, err := io.ReadAll(res.Body)
				assert.NoError(t, err)
				assert.Equal(t, "-modified-next", string(body))

				client.CloseIdleConnections()

				err = s.Shutdown()
				assert.NoError(t, err)

				wg.Wait()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {

			module, err := os.ReadFile(generatedModules[testCase.Module])
			require.NoError(t, err)

			scaleFunc := &scalefunc.ScaleFunc{
				Version:   scalefunc.V1Alpha,
				Name:      "TestName",
				Tag:       "TestTag",
				Signature: "http",
				Language:  scalefunc.Go,
				Function:  module,
			}
			testCase.Run(scaleFunc, t)
		})
	}
}
