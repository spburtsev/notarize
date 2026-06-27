<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import { Field, FieldLabel } from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import Pagination from '$lib/components/pagination.svelte';
	import { enhance } from '$app/forms';
	import { humanizeEnum } from '$lib/fmt.js';
	import type { IssueStatus } from '$lib/api/index.js';

	let { data, form } = $props();
	let createOpen = $state(false);

	function statusBadgeProps(status: IssueStatus) {
		switch (status) {
			case 'OPEN':
				return { variant: 'secondary', class: undefined } as const;
			case 'CLOSED':
				return {
					variant: 'secondary',
					class: 'bg-green-500 text-white dark:bg-green-600'
				} as const;
			default: // IN_PROGRESS
				return { variant: 'default', class: undefined } as const;
		}
	}
</script>

<div class="flex flex-col gap-6 p-6">
	<Card.Root>
		<Card.Header class="flex flex-row items-center justify-between">
			<Card.Title>Issues</Card.Title>
			<Button size="sm" onclick={() => (createOpen = true)}>Create</Button>
		</Card.Header>
		<Card.Content>
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>Title</Table.Head>
						<Table.Head>Status</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each data.issues as i (i.id)}
						<Table.Row>
							<Table.Cell>
								<a class="font-medium hover:underline" href="/issues/{i.id}">{i.title}</a>
							</Table.Cell>
							<Table.Cell
								><Badge {...statusBadgeProps(i.status)}>{humanizeEnum(i.status)}</Badge></Table.Cell
							>
						</Table.Row>
					{:else}
						<Table.Row>
							<Table.Cell colspan={2} class="text-muted-foreground">No issues yet.</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
			<Pagination total={data.total} offset={data.offset} limit={data.limit} />
		</Card.Content>
	</Card.Root>

	<Dialog.Root bind:open={createOpen}>
		<Dialog.Content>
			<Dialog.Header>
				<Dialog.Title>New issue</Dialog.Title>
			</Dialog.Header>
			<form
				method="POST"
				action="?/create"
				use:enhance={() =>
					async ({ result, update }) => {
						await update();
						if (result.type === 'success') createOpen = false;
					}}
				class="flex flex-col gap-4"
			>
				<Field>
					<FieldLabel for="title">Title</FieldLabel>
					<Input id="title" name="title" placeholder="Issue title..." required />
				</Field>
				<Field>
					<FieldLabel for="description">Description (optional)</FieldLabel>
					<Input id="description" name="description" placeholder="Details..." />
				</Field>
				{#if form?.error}<p class="text-destructive text-sm">{form.error}</p>{/if}
				<Dialog.Footer>
					<Button type="submit">Create issue</Button>
				</Dialog.Footer>
			</form>
		</Dialog.Content>
	</Dialog.Root>
</div>
