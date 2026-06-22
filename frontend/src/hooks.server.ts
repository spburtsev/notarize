import { sequence } from '@sveltejs/kit/hooks';
import { error, redirect, type Handle } from '@sveltejs/kit';
import { BACKEND_URL } from '$app/env/private';

const PROXY_PATH = "/api-proxy";

const handleApiProxy: Handle = async ({ event, resolve }) => {
  if (!event.url.pathname.startsWith(PROXY_PATH)) {
    return resolve(event);
  }

  const backendURL = BACKEND_URL;

  const origin = event.request.headers.get("Origin");

  // reject requests that don't come from the webapp, to avoid your proxy being abused.
  if (!origin || new URL(origin).origin !== event.url.origin) {
    error(403, "Request Forbidden.");
  }

  // strip `/api-proxy` from the request path
  const strippedPath = event.url.pathname.substring(PROXY_PATH.length);

  // build the new URL path with your API base URL, the stripped path and the query string
  const urlPath = `${backendURL}${strippedPath}${event.url.search}`;
  const proxiedUrl = new URL(urlPath);

  const headers = new Headers(event.request.headers);
  const session = event.cookies.get("session");
  if (session) {
    headers.set("Authorization", `Bearer ${session}`);
  }

  // Strip off header added by SvelteKit yet forbidden by underlying HTTP request
  // library `undici`.
  // https://github.com/nodejs/undici/issues/1470
  //   headers.delete("connection");

  return fetch(proxiedUrl.toString(), {
    body: event.request.body,
    method: event.request.method,
    headers,
  }).catch((err) => {
    console.log("Could not proxy API request: ", err);
    throw err;
  });
};

const handleAuth: Handle = async ({ event, resolve }) => {
  if (event.route.id?.startsWith("/(protected)") && !event.cookies.get("session")) {
    redirect(303, "/login");
  }

  return resolve(event);
};

export const handle: Handle = sequence(handleApiProxy, handleAuth);