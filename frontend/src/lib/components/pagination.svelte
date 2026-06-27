<script lang="ts">
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button/index.js';

	let { total, offset, limit }: { total: number; offset: number; limit: number } = $props();

	const current = $derived(Math.floor(offset / limit) + 1);
	const pages = $derived(Math.max(1, Math.ceil(total / limit)));
	const hasPrev = $derived(offset > 0);
	const hasNext = $derived(offset + limit < total);

	function href(o: number): string {
		const u = new URL(page.url);
		if (o <= 0) u.searchParams.delete('offset');
		else u.searchParams.set('offset', String(o));
		return u.pathname + u.search;
	}
</script>

<div class="flex items-center justify-between pt-4">
	<span class="text-muted-foreground text-sm">Page {current} of {pages} ({total} total)</span>
	<div class="flex gap-2">
		{#if hasPrev}
			<Button variant="outline" size="sm" href={href(offset - limit)}>Previous</Button>
		{:else}
			<Button variant="outline" size="sm" disabled>Previous</Button>
		{/if}
		{#if hasNext}
			<Button variant="outline" size="sm" href={href(offset + limit)}>Next</Button>
		{:else}
			<Button variant="outline" size="sm" disabled>Next</Button>
		{/if}
	</div>
</div>
