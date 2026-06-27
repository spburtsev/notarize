import { fail, redirect } from '@sveltejs/kit';
import { login } from '$lib/api';
import { serverApi } from '$lib/server/api';
import type { Actions } from './$types';
import { setSession } from '$lib/server/session';

export const actions: Actions = {
	default: async (event) => {
		const { cookies } = event;
		const data = await event.request.formData();
		const email = String(data.get('email') ?? '');
		const password = String(data.get('password') ?? '');

		if (!email || !password) {
			return fail(400, { email, error: 'Email and password are required.' });
		}

		let result, error, response;
		try {
			({
				data: result,
				error,
				response
			} = await login({
				client: serverApi(event),
				body: { email, password }
			}));
		} catch {
			return fail(502, { email, error: 'Could not reach the authentication server.' });
		}

		if (error || !result) {
			if (response?.status === 401) {
				return fail(401, { email, error: 'Invalid email or password.' });
			}
			return fail(502, { email, error: 'Could not reach the authentication server.' });
		}

		setSession(cookies, result.token, result.user);
		redirect(303, '/dashboard');
	}
};
