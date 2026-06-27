import { getDocumentResult } from '$lib/api';
import { propagateErr, serverApi } from '$lib/server/api';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async (event) => {
	const { data } = await propagateErr(getDocumentResult({
		client: serverApi(event),
		path: { documentId: event.params.documentId }
	}));
	return Response.json(data);
};
