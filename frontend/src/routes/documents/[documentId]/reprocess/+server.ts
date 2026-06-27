import { reprocessDocument } from '$lib/api';
import { serverApi, propagateErr } from '$lib/server/api';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async (event) => {
	const { data } = await propagateErr(reprocessDocument({
		client: serverApi(event),
		path: { documentId: event.params.documentId }
	}));
	return Response.json(data);
};
