<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import AppSidebar from '$lib/components/app-sidebar.svelte';
	import SiteHeader from '$lib/components/site-header.svelte';
	import { page } from '$app/state';

	const { children } = $props();
	const breadcrumbs = $derived.by(() => {
		if (page.data.breadcrumbs) return page.data.breadcrumbs;
		const seg = page.url.pathname.split('/').filter(Boolean)[0];
		return seg ? [{ label: seg.charAt(0).toUpperCase() + seg.slice(1), url: '/' + seg }] : [];
	});
</script>

<Sidebar.Provider
	style="--sidebar-width: calc(var(--spacing) * 72); --header-height: calc(var(--spacing) * 12);"
>
	<AppSidebar variant="inset" />
	<Sidebar.Inset>
		<SiteHeader {breadcrumbs} />
		<div class="flex flex-1 flex-col">
			<div class="@container/main flex flex-1 flex-col gap-2">
				{@render children()}
			</div>
		</div>
	</Sidebar.Inset>
</Sidebar.Provider>
