<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import { Field, FieldLabel } from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import Pagination from '$lib/components/pagination.svelte';
	import { enhance } from '$app/forms';

	let { data, form } = $props();
</script>

<div class="flex flex-col gap-6 p-6">
	<Card.Root>
		<Card.Header><Card.Title>New project</Card.Title></Card.Header>
		<Card.Content>
			<form method="POST" action="?/create" use:enhance class="flex items-end gap-3">
				<Field class="flex-1">
					<FieldLabel for="name">Name</FieldLabel>
					<Input id="name" name="name" placeholder="Project name..." required />
				</Field>
				<Button type="submit">Create</Button>
			</form>
			{#if form?.error}<p class="text-destructive mt-2 text-sm">{form.error}</p>{/if}
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header><Card.Title>Projects</Card.Title></Card.Header>
		<Card.Content>
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>Name</Table.Head>
						<Table.Head>Created</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each data.projects as p (p.id)}
						<Table.Row>
							<Table.Cell>
								<a class="font-medium hover:underline" href="/projects/{p.id}">{p.name}</a>
							</Table.Cell>
							<Table.Cell>{new Date(p.created_at).toLocaleDateString()}</Table.Cell>
						</Table.Row>
					{:else}
						<Table.Row>
							<Table.Cell colspan={2} class="text-muted-foreground">No projects yet.</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
			<Pagination total={data.total} offset={data.offset} limit={data.limit} />
		</Card.Content>
	</Card.Root>
</div>
