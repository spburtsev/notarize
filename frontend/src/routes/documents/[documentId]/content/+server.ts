import { error } from '@sveltejs/kit';
import { BACKEND_URL } from '$app/env/private';
import type { RequestHandler } from './$types';
import { getSessionToken } from '$lib/server/session';

export const GET: RequestHandler = async ({ params, cookies, fetch }) => {
	const session = getSessionToken(cookies);
	const res = await fetch(`${BACKEND_URL}/documents/${params.documentId}/content`, {
		headers: session ? { Authorization: `Bearer ${session}` } : undefined
	});
	if (!res.ok) error(res.status, 'Could not load document.');
	return new Response(res.body, {
		headers: { 'Content-Type': 'application/pdf', 'Content-Disposition': 'inline' }
	});
};
