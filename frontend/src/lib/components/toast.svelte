<script lang="ts">
	import { CircleCheck, CircleX, Info, TriangleAlert, X } from '@lucide/svelte';
	import { cn } from '$lib/utils.js';

	export type ToastVariant = 'success' | 'error' | 'info' | 'warning';

	let {
		variant = 'info',
		message,
		description,
		closeToast
	}: {
		variant?: ToastVariant;
		message: string;
		description?: string;
		closeToast?: () => void;
	} = $props();

	const icons = { success: CircleCheck, error: CircleX, info: Info, warning: TriangleAlert };
	const accent = {
		success: 'text-green-600',
		error: 'text-destructive',
		info: 'text-foreground',
		warning: 'text-amber-600'
	};
	const Icon = $derived(icons[variant]);
</script>

<div
	class={cn(
		'bg-popover text-popover-foreground flex w-[356px] max-w-[calc(100vw-2rem)] items-start gap-3 rounded-md border p-3 text-sm shadow-md'
	)}
>
	<Icon class={cn('mt-0.5 size-4 shrink-0', accent[variant])} />
	<div class="min-w-0 flex-1">
		<p class="font-medium">{message}</p>
		{#if description}<p class="text-muted-foreground mt-0.5 text-xs">{description}</p>{/if}
	</div>
	{#if closeToast}
		<button
			type="button"
			onclick={closeToast}
			class="text-muted-foreground hover:text-foreground"
			aria-label="Dismiss"
		>
			<X class="size-4" />
		</button>
	{/if}
</div>
