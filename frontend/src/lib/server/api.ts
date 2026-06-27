import { type RequestEvent, error } from '@sveltejs/kit';
import { BACKEND_URL } from '$app/env/private';
import { createClient, createConfig } from '$lib/api/client';
import { getSessionToken } from './session';

export function serverApi(event: Pick<RequestEvent, 'fetch' | 'cookies'>) {
	const session = getSessionToken(event.cookies);
	return createClient(
		createConfig({
			baseUrl: BACKEND_URL,
			fetch: event.fetch,
			headers: session ? { Authorization: `Bearer ${session}` } : undefined
		})
	);
}

export async function propagateErr<T extends { data: unknown; error: unknown }>(
	pending: Promise<T & { request?: Request; response?: Response }>
): Promise<{ data: NonNullable<T['data']>; request: Request; response: Response }> {
	const { data, error: callError, request, response } = await pending;
	if (callError) {
		const message =
			typeof callError === 'object' && callError !== null && 'message' in callError
				? String((callError as { message: unknown }).message)
				: 'An error occurred while processing the request.';
		error(response?.status ?? 500, message);
	}
	return {
		data: data as NonNullable<T['data']>,
		request: request as Request,
		response: response as Response
	};
}
