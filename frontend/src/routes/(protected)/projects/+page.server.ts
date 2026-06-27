import { fail } from '@sveltejs/kit';
import { listProjects, createProject } from '$lib/api';
import { serverApi } from '$lib/server/api';
import { PAGE_SIZE, parseOffset } from '$lib/pagination';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	const offset = parseOffset(event.url);
	const { data } = await listProjects({
		client: serverApi(event),
		query: { limit: PAGE_SIZE, offset }
	});
	return { projects: data?.items ?? [], total: data?.total ?? 0, offset, limit: PAGE_SIZE };
};

export const actions: Actions = {
	create: async (event) => {
		const fd = await event.request.formData();
		const name = String(fd.get('name') ?? '').trim();
		if (!name) return fail(400, { error: 'Name is required.' });
		const { error } = await createProject({ client: serverApi(event), body: { name } });
		if (error) return fail(502, { error: 'Could not create project.' });
		return { success: true };
	}
};
