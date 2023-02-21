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

import {Context, HttpStringList} from "@loopholelabs/scale-signature-http";

import { Runtime } from "@loopholelabs/scale";

import {NextRequest, NextResponse} from 'next/server';

// https://vercel.com/docs/concepts/functions/edge-functions#creating-edge-functions

export class NextJS {
  private _runtime: Runtime<Context>;
  constructor(runtime: Runtime<Context>) {
    this._runtime = runtime;
  }

  Handler() {
    return async (req: NextRequest) => {
      const i = await this._runtime.Instance(null);
      await NextJS.fromRequest(i.Context(), req);
      i.Run();
      return NextJS.toResponse(i.Context());
    };
  }

  static async fromRequest(ctx: Context, req: NextRequest): Promise<Context> {
    ctx.Request.Protocol = "HTTP/1.1"
    ctx.Request.Method = req.method;
    ctx.Request.IP = req.headers.get("x-forwarded-for") || req.headers.get("x-real-ip") || req.ip || "";
    ctx.Request.ContentLength = BigInt(req.headers.get("content-length") || 0);

    ctx.Request.Body = new Uint8Array(0);
    if (req.body!=null) {
      ctx.Request.Body = new Uint8Array(await (await req.blob()).arrayBuffer());
    }

    if (ctx.Request.ContentLength < 1) {
      ctx.Request.ContentLength = BigInt(ctx.Request.Body.length);
    }

    for(let k in req.headers) {
      let vals = req.headers.get(k);
      let sl: string[] = [];
      if (vals !== null) {
        sl.push(vals);
      }
      const v = new HttpStringList(sl);
      ctx.Request.Headers.set(k, v);
    }

    return ctx;
  }

  static toResponse(ctx: Context): NextResponse {
    const headers = new Headers();

    ctx.Response.Headers.forEach((v, k) => {
      let vals = v.Value;
      if (vals !== undefined) {
        for(let v of vals.values()) {
          headers.set(k, v);
        }
      }
    })

    return new NextResponse(ctx.Response.Body, {
      status: ctx.Response.StatusCode || 200,
      headers: headers
    });
  }
}