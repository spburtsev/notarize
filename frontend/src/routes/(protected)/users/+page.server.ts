import { fail, redirect } from '@sveltejs/kit';
import { listUsers, createUser } from '$lib/api';
import type { UserRole } from '$lib/api';
import { serverApi } from '$lib/server/api';
import { PAGE_SIZE, parseOffset } from '$lib/pagination';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	const { user } = await event.parent();
	if (user?.role !== 'ADMIN') redirect(303, '/');

	const offset = parseOffset(event.url);
	const { data } = await listUsers({
		client: serverApi(event),
		query: { limit: PAGE_SIZE, offset }
	});
	return { users: data?.items ?? [], total: data?.total ?? 0, offset, limit: PAGE_SIZE };
};

export const actions: Actions = {
	create: async (event) => {
		const fd = await event.request.formData();
		const email = String(fd.get('email') ?? '').trim();
		const first_name = String(fd.get('first_name') ?? '').trim();
		const last_name = String(fd.get('last_name') ?? '').trim();
		const password = String(fd.get('password') ?? '');
		const role = String(fd.get('role') ?? 'EMPLOYEE') as UserRole;
		if (!email || !first_name || !last_name || !password) {
			return fail(400, { error: 'All fields are required.' });
		}
		const { error, response } = await createUser({
			client: serverApi(event),
			body: { email, first_name, last_name, password, role }
		});
		if (error) {
			if (response?.status === 403)
				return fail(403, { error: 'Admin role required to create users.' });
			return fail(502, { error: 'Could not create user (email may already exist).' });
		}
		return { success: true };
	}
};
