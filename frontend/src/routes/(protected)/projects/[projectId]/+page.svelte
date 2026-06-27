<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import { Field, FieldLabel } from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import Pagination from '$lib/components/pagination.svelte';
	import { enhance } from '$app/forms';

	let { data, form } = $props();
</script>

<div class="flex flex-col gap-6 p-6">
	<h1 class="text-xl font-semibold">{data.project?.name ?? 'Project'}</h1>

	<Card.Root>
		<Card.Header><Card.Title>New issue</Card.Title></Card.Header>
		<Card.Content>
			<form method="POST" action="?/create" use:enhance class="flex flex-col gap-3">
				<Field>
					<FieldLabel for="title">Title</FieldLabel>
					<Input id="title" name="title" placeholder="Issue title..." required />
				</Field>
				<Field>
					<FieldLabel for="description">Description (optional)</FieldLabel>
					<Input id="description" name="description" placeholder="Details..." />
				</Field>
				<div><Button type="submit">Create issue</Button></div>
			</form>
			{#if form?.error}<p class="text-destructive mt-2 text-sm">{form.error}</p>{/if}
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header><Card.Title>Issues</Card.Title></Card.Header>
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
							<Table.Cell><Badge variant="secondary">{i.status}</Badge></Table.Cell>
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
</div>
