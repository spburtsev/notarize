<script lang="ts">
	import * as Command from '$lib/components/ui/command/index.js';
	import { SquarePen } from '@lucide/svelte';
	import { FolderIcon } from '@lucide/svelte';
	import { UserIcon } from '@lucide/svelte';

	let { open = $bindable(false) }: { open?: boolean } = $props();
</script>

<Command.Dialog bind:open>
	<Command.Input placeholder="Type an item type or search..." />
	<Command.List>
		<Command.Empty>No results found.</Command.Empty>
		<Command.Group heading="Suggestions">
			<Command.Item>
				<FolderIcon class="me-2 size-4" />
				<span>Project</span>
			</Command.Item>
			<Command.Item>
				<SquarePen class="me-2 size-4" />
				<span>Issue</span>
			</Command.Item>
			<Command.Item>
				<UserIcon class="me-2 size-4" />
				<span>User</span>
			</Command.Item>
		</Command.Group>
		<Command.Separator />
	</Command.List>
</Command.Dialog>
