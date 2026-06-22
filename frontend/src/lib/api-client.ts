import { createClient, createConfig } from '$lib/api/client';

export const api = createClient(createConfig({ baseUrl: '/api-proxy' }));
