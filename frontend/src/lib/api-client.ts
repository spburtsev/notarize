import { createClient, createConfig } from '$lib/api/client';

export const clientApi = createClient(createConfig({ baseUrl: '/api-proxy' }));
