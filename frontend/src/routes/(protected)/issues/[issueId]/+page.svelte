<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Pencil } from '@lucide/svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { enhance } from '$app/forms';
	import DocWorkbench from './doc-workbench.svelte';

	let { data, form: uploadForm } = $props();

	const selectedDocId = $derived(page.url.searchParams.get('document'));
	const uploadError = $derived(uploadForm?.error);

	function selectDoc(id: string) {
		const url = new URL(page.url);
		url.searchParams.set('document', id);
		goto(url, { keepFocus: true, noScroll: true });
	}

	const description = $derived(
		data.issue === null ? '...' : (data.issue.description ?? 'No description available.')
	);

	const canEdit = $derived(data.user?.role === 'ADMIN' || data.user?.role === 'MANAGER');
	let editing = $state(false);
	let draft = $state('');
	function startEdit() {
		draft = data.issue?.description ?? '';
		editing = true;
	}
</script>

<div class="flex flex-col gap-6 p-6">
	<Tabs.Root value="details">
		<Tabs.List>
			<Tabs.Trigger value="details">Details</Tabs.Trigger>
			<Tabs.Trigger value="members">Members</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content value="details">
			{#if editing}
				<form
					method="POST"
					action="?/update"
					use:enhance={() =>
						async ({ result, update }) => {
							await update();
							if (result.type === 'success') editing = false;
						}}
					class="flex flex-col gap-2"
				>
					<textarea
						name="description"
						bind:value={draft}
						rows="2"
						placeholder="Issue description..."
						class="border-input border bg-muted px-3 py-2 text-sm"></textarea>
					{#if uploadForm?.error}<p class="text-destructive text-sm">{uploadForm.error}</p>{/if}
					<div class="flex gap-2">
						<Button type="submit" size="sm">Save</Button>
						<Button type="button" size="sm" variant="ghost" onclick={() => (editing = false)}>
							Cancel
						</Button>
					</div>
				</form>
			{:else}
				<div class="flex items-start gap-2">
					<p class="text-sm flex-1 bg-muted border px-3 py-2">{description}</p>
					{#if canEdit}
						<Button
							size="icon-sm"
							variant="ghost"
							onclick={startEdit}
							aria-label="Edit description"
						>
							<Pencil />
						</Button>
					{/if}
				</div>
			{/if}
		</Tabs.Content>
		<Tabs.Content value="members">
          <h2>Employees:</h2>
          <h2>Invitees:</h2>
        </Tabs.Content>
	</Tabs.Root>

	<Card.Root>
		<Card.Content>
			<DocWorkbench {selectedDocId} onSelect={selectDoc} formError={uploadError} {data} />
		</Card.Content>
	</Card.Root>
</div>
