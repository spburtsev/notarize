<script lang="ts">
	import ChartBarIcon from '@tabler/icons-svelte/icons/chart-bar';
	import DashboardIcon from '@tabler/icons-svelte/icons/dashboard';
	import FolderIcon from '@tabler/icons-svelte/icons/folder';
	import InnerShadowTopIcon from '@tabler/icons-svelte/icons/inner-shadow-top';
	import UsersIcon from '@tabler/icons-svelte/icons/users';
	import NavMain from './nav-main.svelte';
	import NavUser from './nav-user.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';
	import type { User } from '$lib/api';

	let { user, ...restProps }: { user?: User | null } & ComponentProps<typeof Sidebar.Root> =
		$props();

	const baseNav = [
		{ title: 'Dashboard', url: '/dashboard', icon: DashboardIcon },
		{ title: 'Analytics', url: '/analytics', icon: ChartBarIcon },
		{ title: 'Projects', url: '/projects', icon: FolderIcon },
		{ title: 'Users', url: '/users', icon: UsersIcon }
	];
	// Users management is admin-only.
	const navMain = $derived(
		user?.role === 'ADMIN' ? baseNav : baseNav.filter((item) => item.url !== '/users')
	);
</script>

<Sidebar.Root collapsible="icon" {...restProps}>
	<Sidebar.Header>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton class="data-[slot=sidebar-menu-button]:p-1.5!">
					{#snippet child({ props })}
						<a href="##" {...props}>
							<InnerShadowTopIcon class="size-5!" />
							<span class="text-base font-semibold">Notarize</span>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Header>
	<Sidebar.Content>
		<NavMain items={navMain} />
	</Sidebar.Content>
	<Sidebar.Footer>
		<NavUser {user} />
	</Sidebar.Footer>
</Sidebar.Root>
