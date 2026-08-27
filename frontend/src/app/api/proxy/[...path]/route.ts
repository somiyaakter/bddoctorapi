import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_BASE_URL;
const API_KEY = process.env.API_KEY;

export async function GET(req: NextRequest, { params }: { params: Promise<{ path: string[] }> }) {
  const { path } = await params;
  const targetUrl = `${API_BASE_URL}/api/v1/${path.join("/")}${req.nextUrl.search}`;

  const res = await fetch(targetUrl, {
    headers: { "X-API-Key": API_KEY ?? "" },
  });

  const body = await res.json();
  return NextResponse.json(body, { status: res.status });
}