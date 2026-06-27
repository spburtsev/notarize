import { fail } from '@sveltejs/kit';
import { getIssue, getProject, listDocuments, createDocument, updateIssue } from '$lib/api';
import { serverApi } from '$lib/server/api';
import type { Actions, PageServerLoad } from './$types';
import { propagateErr } from '$lib/server/api';

export const load: PageServerLoad = async (event) => {
	const client = serverApi(event);
	const [issue, docs] = await Promise.all([
		propagateErr(getIssue({ client, path: { issueId: event.params.issueId } })),
		propagateErr(listDocuments({ client, query: { issueId: event.params.issueId, limit: 200 } }))
	]);

	const { data: project } = await getProject({
		client,
		path: { projectId: issue.data.project_id }
	});

	const breadcrumbs = [
		{ label: 'Projects', url: '/projects' },
		...(project ? [{ label: project.name, url: `/projects/${project.id}` }] : []),
		{ label: issue?.data.title ?? 'Issue', url: `/issues/${event.params.issueId}` }
	];

	return { issue: issue.data, documents: docs.data?.items ?? [], breadcrumbs };
};

export const actions: Actions = {
	upload: async (event) => {
		const fd = await event.request.formData();
		const file = fd.get('file');
		if (!(file instanceof File) || file.size === 0) {
			return fail(400, { error: 'Choose a file to upload.' });
		}
		const { error, response } = await createDocument({
			client: serverApi(event),
			body: { file, issue_id: event.params.issueId }
		});
		if (error) {
			if (response?.status === 400) return fail(400, { error: 'Only PDF documents are accepted.' });
			return fail(502, { error: 'Upload failed.' });
		}
		return { success: true };
	},
	update: async (event) => {
		const fd = await event.request.formData();
		const description = String(fd.get('description') ?? '');
		const { error, response } = await updateIssue({
			client: serverApi(event),
			path: { issueId: event.params.issueId },
			body: { description: description || null }
		});
		if (error) {
			if (response?.status === 403) {
				return fail(403, { error: 'Only managers and admins can edit issues.' });
            }
			return fail(502, { error: `Could not update issue: ${error.message}.` });
		}
		return { success: true };
	}
};
