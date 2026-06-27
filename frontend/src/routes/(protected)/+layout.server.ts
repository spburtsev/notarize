import { requireSession } from '$lib/server/session';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = ({ cookies }) => {
	const { user } = requireSession(cookies);
	return { user };
};
