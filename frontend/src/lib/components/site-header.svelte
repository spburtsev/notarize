<script lang="ts">
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';

	let { breadcrumbs }: { breadcrumbs: { label: string; url: string }[] } = $props();

	const pageTitle = $derived(breadcrumbs[breadcrumbs.length - 1]?.label ?? '');
</script>

<header
	class="flex h-(--header-height) shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-(--header-height)"
>
	<div class="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
		<Sidebar.Trigger class="-ms-1" />
		<Separator orientation="vertical" class="mx-2 data-[orientation=vertical]:h-4" />
		{#if breadcrumbs.length > 1}
			<Breadcrumb.Root>
				<Breadcrumb.List>
					{#each breadcrumbs as crumb, i (crumb.label)}
						<Breadcrumb.Item>
							{#if i < breadcrumbs.length - 1}
								<Breadcrumb.Link href={crumb.url}>{crumb.label}</Breadcrumb.Link>
							{:else}
								<Breadcrumb.Page>{crumb.label}</Breadcrumb.Page>
							{/if}
						</Breadcrumb.Item>
						{#if i < breadcrumbs.length - 1}
							<Breadcrumb.Separator />
						{/if}
					{/each}
				</Breadcrumb.List>
			</Breadcrumb.Root>
		{:else}
			<h1 class="text-base font-medium">{pageTitle}</h1>
		{/if}
	</div>
</header>
