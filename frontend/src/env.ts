import { defineEnvVars } from '@sveltejs/kit/hooks';

export const variables = defineEnvVars({
	BACKEND_URL: {}
});
