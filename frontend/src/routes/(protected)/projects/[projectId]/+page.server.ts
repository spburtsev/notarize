import { fail } from '@sveltejs/kit';
import { getProject, listIssues, createIssue } from '$lib/api';
import { serverApi } from '$lib/server/api';
import { PAGE_SIZE, parseOffset } from '$lib/pagination';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	const client = serverApi(event);
	const offset = parseOffset(event.url);
	const [{ data: project }, { data: issues }] = await Promise.all([
		getProject({ client, path: { projectId: event.params.projectId } }),
		listIssues({ client, query: { projectId: event.params.projectId, limit: PAGE_SIZE, offset } })
	]);
	return {
		project: project ?? null,
		issues: issues?.items ?? [],
		total: issues?.total ?? 0,
		offset,
		limit: PAGE_SIZE
	};
};

export const actions: Actions = {
	create: async (event) => {
		const fd = await event.request.formData();
		const title = String(fd.get('title') ?? '').trim();
		const description = String(fd.get('description') ?? '').trim();
		if (!title) return fail(400, { error: 'Title is required.' });
		const { error } = await createIssue({
			client: serverApi(event),
			body: { project_id: event.params.projectId, title, description: description || null }
		});
		if (error) return fail(502, { error: 'Could not create issue.' });
		return { success: true };
	}
};
