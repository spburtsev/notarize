<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import { Field, FieldLabel } from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import Pagination from '$lib/components/pagination.svelte';
	import { enhance } from '$app/forms';
	import { humanizeEnum } from '$lib/fmt.js';

	let { data, form } = $props();

	const roles = ['ADMIN', 'MANAGER', 'EMPLOYEE', 'INVITEE'];
	let role = $state('EMPLOYEE');
</script>

<div class="flex flex-col gap-6 p-6">
	<Card.Root>
		<Card.Header><Card.Title>New user</Card.Title></Card.Header>
		<Card.Content>
			<form method="POST" action="?/create" use:enhance class="flex flex-col gap-3">
				<div class="grid grid-cols-2 gap-3">
					<Field>
						<FieldLabel for="first_name">First name</FieldLabel>
						<Input id="first_name" name="first_name" required />
					</Field>
					<Field>
						<FieldLabel for="last_name">Last name</FieldLabel>
						<Input id="last_name" name="last_name" required />
					</Field>
				</div>
				<Field>
					<FieldLabel for="email">Email</FieldLabel>
					<Input id="email" name="email" type="email" placeholder="user@example.com" required />
				</Field>
				<div class="grid grid-cols-2 gap-3">
					<Field>
						<FieldLabel for="password">Password</FieldLabel>
						<Input id="password" name="password" type="password" required />
					</Field>
					<Field>
						<FieldLabel for="role">Role</FieldLabel>
						<Select.Root type="single" name="role" bind:value={role}>
							<Select.Trigger id="role" class="w-full">{humanizeEnum(role)}</Select.Trigger>
							<Select.Content>
								{#each roles as r (r)}
									<Select.Item value={r} label={humanizeEnum(r)} />
								{/each}
							</Select.Content>
						</Select.Root>
					</Field>
				</div>
				<div><Button type="submit">Create user</Button></div>
			</form>
			{#if form?.error}<p class="text-destructive mt-2 text-sm">{form.error}</p>{/if}
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header><Card.Title>Users</Card.Title></Card.Header>
		<Card.Content>
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>Email</Table.Head>
						<Table.Head>Name</Table.Head>
						<Table.Head>Role</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each data.users as u (u.id)}
						<Table.Row>
							<Table.Cell class="font-medium">{u.email}</Table.Cell>
							<Table.Cell>{u.first_name} {u.last_name}</Table.Cell>
							<Table.Cell><Badge variant="secondary">{humanizeEnum(u.role)}</Badge></Table.Cell>
						</Table.Row>
					{:else}
						<Table.Row>
							<Table.Cell colspan={3} class="text-muted-foreground">
								No users visible (admin only).
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
			<Pagination total={data.total} offset={data.offset} limit={data.limit} />
		</Card.Content>
	</Card.Root>
</div>
