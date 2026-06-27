import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { eraseSession } from '$lib/server/session';

export const POST: RequestHandler = ({ cookies }) => {
	eraseSession(cookies);
	redirect(303, '/login');
};
