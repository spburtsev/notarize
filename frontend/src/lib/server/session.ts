import type { User } from '$lib/api';
import { type Cookies, redirect } from '@sveltejs/kit';

const SESSION_KEY = 'session';
const USER_KEY = 'user';
const SESSION_MAX_AGE = 60 * 60 * 24;

export function getSessionToken(cookies: Cookies) {
	return cookies.get(SESSION_KEY);
}

export function requireSession(cookies: Cookies) {
	const session = cookies.get(SESSION_KEY);
	if (!session) {
		redirect(303, '/login');
	}
	const userString = cookies.get(USER_KEY);
	if (!userString) {
		cookies.delete(SESSION_KEY, { path: '/' });
		redirect(303, '/login');
	}

	const user = JSON.parse(userString) as User;
	// TODO: Validate the user object structure if necessary

	return { session, user };
}

export function eraseSession(cookies: Cookies) {
	cookies.delete(SESSION_KEY, { path: '/' });
	cookies.delete(USER_KEY, { path: '/' });
}

export function setSession(cookies: Cookies, token: string, user: User) {
	const cookieOpts = {
		httpOnly: true,
		sameSite: 'lax',
		path: '/',
		maxAge: SESSION_MAX_AGE
	} as const;
	cookies.set('session', token, cookieOpts);
	cookies.set('user', JSON.stringify(user), cookieOpts);
}
