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

import {
  Headers,
  Request,
  Response,
} from 'node-fetch';

if (!global.fetch) {
//  (global as any).fetch = fetch;
  (global as any).Headers = Headers;
  (global as any).Request = Request;
  (global as any).Response = Response;
}

import { TextEncoder, TextDecoder } from "util";

window.TextEncoder = TextEncoder;
window.TextDecoder = TextDecoder as typeof window["TextDecoder"];

import * as fs from "fs";

import { NextJS } from "./index";
import { NextRequest } from 'next/server';
import { New } from "@loopholelabs/scale";
import { ScaleFunc, V1Alpha, Go } from "@loopholelabs/scalefile";
import * as signature from "@loopholelabs/scale-signature-http";

describe("nextAdapter", () => {

  it("Can convert Request to Context", async () => {
    const bodyData = '{"foo": "bar"}';
    const request = new NextRequest('https://example.com', {method: 'POST', body: bodyData});

    let ctx = signature.New();
    ctx = await NextJS.fromRequest(ctx, request);

    if (request.body != null ) {
      expect(ctx.Request.Method).toBe(request.method);
      expect(ctx.Request.Protocol).toBe("HTTP/1.1");
      expect(Number(ctx.Request.ContentLength)).toBe(bodyData.length);
      const reqBody = new TextDecoder().decode(ctx.Request.Body);
      expect(reqBody).toBe(bodyData);
    }
  });

  it("Can convert Context to Response", async () => {
    const body = new TextEncoder().encode("Hello world");
    const headers = new Map<string, signature.HttpStringList>;
    headers.set("MIDDLEWARE", new signature.HttpStringList(["Hello"]));
    const c = new signature.Context();
    c.Response.Body = body;
    c.Response.Headers = headers;

    const response = NextJS.toResponse(c);

    let b = await (await response.blob()).arrayBuffer();
    const outbodybytes = new Uint8Array(b);
    const outbody = new TextDecoder().decode(outbodybytes);

    expect(outbody).toBe("Hello world");
    expect(response.status).toBe(200);

    // Check for the header
    const hkey = response.headers.get("MIDDLEWARE");
    expect(hkey).toBe("Hello");
  });

  it("Can run end-to-end", async () => {
    const modNext = fs.readFileSync("tests/modules/next/next.wasm");

    const fn = new ScaleFunc(V1Alpha, "Test.Next", "Test.Tag", "ExampleName@ExampleVersion", Go, modNext);
    const r = await New([fn]);
    const adapter = new NextJS(r);
    const handler = adapter.Handler();

    const bodyData = '{"foo": "bar"}';
    const request = new NextRequest('https://example.com', {method: 'POST', body: bodyData, headers: {
        "test-x": "test",
      }});
    const res = await handler(request);
    let b = await (await res.blob()).arrayBuffer();
    const outbodybytes = new Uint8Array(b);
    const outbody = new TextDecoder().decode(outbodybytes);

    expect(res.status).toEqual(200);
    expect(outbody).toBe("-modified");
  });

});
