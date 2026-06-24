<script lang="ts">
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Card from "$lib/components/ui/card/index.js";
	import {
		FieldGroup,
		Field,
		FieldLabel,
		FieldDescription,
	} from "$lib/components/ui/field/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { cn } from "$lib/utils.js";
	import { enhance } from "$app/forms";
	import type { HTMLAttributes } from "svelte/elements";

	type LoginFormResult = { email?: string; error?: string };

	let {
		class: className,
		form = null,
		...restProps
	}: HTMLAttributes<HTMLDivElement> & { form?: LoginFormResult | null } = $props();

	const id = $props.id();
</script>

<div class={cn("flex flex-col gap-6", className)} {...restProps}>
	<Card.Root>
		<Card.Header class="text-center">
			<Card.Title class="text-xl">Welcome back</Card.Title>
			<Card.Description>Login with your email and password</Card.Description>
		</Card.Header>
		<Card.Content>
			<form method="POST" use:enhance>
				<FieldGroup>
					<Field>
						<FieldLabel for="email-{id}">Email</FieldLabel>
						<Input
							id="email-{id}"
							name="email"
							type="email"
							placeholder="m@example.com"
							value={form?.email ?? ""}
							required
						/>
					</Field>
					<Field>
						<div class="flex items-center">
							<FieldLabel for="password-{id}">Password</FieldLabel>
							<a href="##" class="ms-auto text-sm underline-offset-4 hover:underline">
								Forgot your password?
							</a>
						</div>
						<Input id="password-{id}" name="password" type="password" required />
					</Field>
					{#if form?.error}
						<p class="text-destructive text-center text-sm">{form.error}</p>
					{/if}
					<Field>
						<Button type="submit">Login</Button>
						<FieldDescription class="text-center">
							Don't have an account? <a href="##">Sign up</a>
						</FieldDescription>
					</Field>
				</FieldGroup>
			</form>
		</Card.Content>
	</Card.Root>
</div>
