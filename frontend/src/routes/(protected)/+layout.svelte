<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import AppSidebar from '$lib/components/app-sidebar.svelte';
	import SiteHeader from '$lib/components/site-header.svelte';
	import { page } from '$app/state';

	const { children } = $props();
	const breadcrumbs = $derived.by(() => {
		switch (page.url.pathname) {
			case '/dashboard':
				return [{ label: 'Dashboard', url: '/dashboard' }];
			case '/analytics':
				return [{ label: 'Analytics', url: '/analytics' }];
			case '/projects':
				return [{ label: 'Projects', url: '/projects' }];
			case '/team':
				return [{ label: 'Team', url: '/team' }];
			default:
				return [];
		}
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
