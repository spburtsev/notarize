import type { RequestEvent } from '@sveltejs/kit';
import { BACKEND_URL } from '$app/env/private';
import { createClient, createConfig } from '$lib/api/client';

export function serverApi(event: Pick<RequestEvent, 'fetch' | 'cookies'>) {
	const session = event.cookies.get('session');
	return createClient(
		createConfig({
			baseUrl: BACKEND_URL,
			fetch: event.fetch,
			headers: session ? { Authorization: `Bearer ${session}` } : undefined
		})
	);
}
