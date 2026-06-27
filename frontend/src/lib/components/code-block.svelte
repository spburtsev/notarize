<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { Copy, Check } from '@lucide/svelte';
	import { cn } from '$lib/utils.js';

	let {
		content,
		label,
		class: className
	}: { content: string | null | undefined; label?: string; class?: string } = $props();

	let copied = $state(false);
	let timer: ReturnType<typeof setTimeout>;

	async function copy() {
		if (!navigator.clipboard) return;
		await navigator.clipboard.writeText(content ?? '');
		copied = true;
		clearTimeout(timer);
		timer = setTimeout(() => (copied = false), 1500);
	}
</script>

<div class={cn('bg-muted overflow-hidden rounded-md border', className)}>
	<div class="flex items-center justify-between gap-2 border-b px-3 py-1.5">
		<span class="text-muted-foreground text-xs font-medium">{label}</span>
		<Button variant="outline" size="sm" onclick={copy} aria-label="Copy">
			{#if copied}<Check /> Copied{:else}<Copy /> Copy{/if}
		</Button>
	</div>
	<pre class="max-h-[70vh] overflow-auto p-3 text-xs whitespace-pre-wrap">{content ?? ''}</pre>
</div>
