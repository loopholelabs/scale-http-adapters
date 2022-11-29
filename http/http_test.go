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
	"context"
	signature "github.com/loopholelabs/scale-signature-http"
	runtime "github.com/loopholelabs/scale/go"
	"github.com/loopholelabs/scale/go/tests/harness"
	"github.com/loopholelabs/scalefile/scalefunc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type TestCase struct {
	Name   string
	Module *harness.Module
	Run    func(*scalefunc.ScaleFunc, *testing.T)
}

func TestHTTP(t *testing.T) {
	nextModule := &harness.Module{
		Name:      "next",
		Path:      "tests/modules/next/next.go",
		Signature: "github.com/loopholelabs/scale-signature-http",
	}

	modules := []*harness.Module{nextModule}

	generatedModules := harness.Setup(t, modules, "github.com/loopholelabs/scale-http-adapters/http/tests/modules")

	var testCases = []TestCase{
		{
			Name:   "Passthrough",
			Module: nextModule,
			Run: func(scaleFunc *scalefunc.ScaleFunc, t *testing.T) {
				r, err := runtime.New(context.Background(), signature.New(), []*scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				adapter := New(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					require.Equal(t, "GET", req.Method)
					body, err := io.ReadAll(req.Body)
					require.NoError(t, err)
					assert.Equal(t, "Test Data", string(body))
					_, err = w.Write([]byte("Hello World"))
					require.NoError(t, err)
				}), r)

				server := httptest.NewServer(adapter)
				defer server.Close()

				req, err := http.NewRequest("GET", server.URL, strings.NewReader("Test Data"))
				assert.NoError(t, err)

				res, err := http.DefaultClient.Do(req)
				assert.NoError(t, err)

				body, err := io.ReadAll(res.Body)
				assert.NoError(t, err)
				assert.Equal(t, "Hello World", string(body))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {

			module, err := os.ReadFile(generatedModules[testCase.Module])
			require.NoError(t, err)

			scaleFunc := &scalefunc.ScaleFunc{
				Version:   "TestVersion",
				Name:      "TestName",
				Signature: "http@v0.1.1",
				Language:  "go",
				Function:  module,
			}
			testCase.Run(scaleFunc, t)
		})
	}
}
