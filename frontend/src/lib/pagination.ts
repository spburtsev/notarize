export const PAGE_SIZE = 25;

export const parseOffset = (url: URL): number =>
	Math.max(0, Number(url.searchParams.get('offset')) || 0);
