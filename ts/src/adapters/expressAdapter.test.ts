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

import { TextEncoder, TextDecoder } from "util";
import express from "express";
import bodyParser from "body-parser";
import * as fs from "fs";
import request from "supertest";

import { ScaleFunc, V1Alpha, Go } from "@loopholelabs/scalefile";
import { HttpContextFactory } from "@loopholelabs/scale-signature-http";
import { GetRuntime } from "@loopholelabs/scale-ts";

import { ExpressAdapter } from "./expressAdapter";

window.TextEncoder = TextEncoder;
window.TextDecoder = TextDecoder as typeof window["TextDecoder"];


describe("expressAdapter", () => {
  const app = express();

  it("Can run a simple e2e", async () => {
    const modHttpEndpoint = fs.readFileSync(
      "./example_modules/http-endpoint.wasm"
    );
    const modHttpMiddleware = fs.readFileSync(
      "./example_modules/http-middleware.wasm"
    );

    const scalefnEndpoint = new ScaleFunc(V1Alpha, "Test.HttpEndpoint", "ExampleName@ExampleVersion", Go, [], modHttpEndpoint);
    const scalefnMiddle = new ScaleFunc(V1Alpha, "Test.HttpMiddleware", "ExampleName@ExampleVersion", Go, [], modHttpMiddleware);

    const signatureFactory = HttpContextFactory;

    const r = await GetRuntime(signatureFactory, [scalefnMiddle, scalefnEndpoint]);

    const adapter = new ExpressAdapter(r);

    app.use(
      bodyParser.raw({
        type: () => true,
      })
    );

    app.use(adapter.getHandler());

    const res = await request(app).post("/blah").send("HELLO WORLD");

    // Make sure everything worked as expected.
    expect(res.statusCode).toEqual(200);
    expect(res.text).toBe("HELLO WORLD");
    expect(res.headers.middleware).toBe("TRUE");
  });
});
