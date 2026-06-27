<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Field, FieldLabel } from '$lib/components/ui/field/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { ButtonGroup } from '$lib/components/ui/button-group/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Item from '$lib/components/ui/item/index.js';
	import CodeBlock from '$lib/components/code-block.svelte';
	import { Upload, FileText, TriangleAlert, RotateCcw } from '@lucide/svelte';
	import { enhance } from '$app/forms';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { humanizeEnum, kb } from '$lib/fmt.js';

	let { data, form } = $props();

	const selected = $derived(page.url.searchParams.get('document'));
	const selectedDoc = $derived(data.documents.find((d) => d.id === selected) ?? null);

	function select(id: string) {
		const url = new URL(page.url);
		url.searchParams.set('document', id);
		goto(url, { keepFocus: true, noScroll: true });
	}

	type Result = {
		status: 'QUEUED' | 'PROCESSING' | 'DONE' | 'FAILED';
		error?: string | null;
		output?: string | null;
	};
	let result = $state<Result | null>(null);
	let dialogOpen = $state(false);
	let uploadOpen = $state(false);
	let restartTrigger = $state(0);

	$effect(() => {
		const id = selected;
		void restartTrigger;
		result = null;
		if (!id) return;
		let stopped = false;
		let timer: ReturnType<typeof setTimeout>;
		const poll = async () => {
			const res = await fetch(`/documents/${id}/result`);
			if (stopped || !res.ok) return;
			result = (await res.json()) as Result;
			if (result?.status === 'QUEUED' || result?.status === 'PROCESSING') {
				timer = setTimeout(poll, 2000);
			}
		};
		poll();
		return () => {
			stopped = true;
			clearTimeout(timer);
		};
	});

	const statusLabel = $derived(
		{
			QUEUED: 'Queued for parsing…',
			PROCESSING: 'Parsing…',
			DONE: 'Parsed',
			FAILED: 'Parsing failed'
		}[result?.status ?? 'QUEUED']
	);
	const dotClass = $derived(
		{
			QUEUED: 'bg-muted-foreground',
			PROCESSING: 'bg-blue-500 animate-pulse',
			DONE: 'bg-green-500',
			FAILED: 'bg-destructive'
		}[result?.status ?? 'QUEUED']
	);

	async function restart() {
		const id = selected;
		if (!id) return;
		await fetch(`/documents/${id}/reprocess`, { method: 'POST' });
		++restartTrigger;
	}

	const description = $derived(
		data.issue === null ? '...' : (data.issue.description ?? 'No description available.')
	);
</script>

<div class="flex flex-col gap-6 p-6">
	<p class="text-md">{description}</p>

	<Card.Root>
		<Card.Content>
			<div class="mb-3 flex items-center justify-between gap-2 text-sm">
				<div class="flex items-center gap-2">
					{#if selectedDoc && result}
						<span class="size-2 rounded-full {dotClass}"></span>
						<span>{statusLabel}</span>
						{#if result.status === 'DONE'}
							<Button size="sm" variant="outline" class="ml-2" onclick={() => (dialogOpen = true)}>
								<FileText />
								View results
							</Button>
						{:else if result.status === 'FAILED'}
							<ButtonGroup class="ml-2">
								{#if result.error}
									<Button size="sm" variant="outline" onclick={() => (dialogOpen = true)}>
										<TriangleAlert />
										Show error
									</Button>
								{/if}
								<Button size="sm" variant="outline" onclick={restart}>
									<RotateCcw />
									Restart
								</Button>
							</ButtonGroup>
						{/if}
					{:else if selectedDoc}
						<Spinner />
						<span class="text-muted-foreground">Checking status...</span>
					{/if}
				</div>
				<Button size="sm" onclick={() => (uploadOpen = true)}>
					<Upload />
					Upload
				</Button>
			</div>
			<div class="flex h-[70vh] gap-4">
				<div class="w-64 shrink-0 space-y-1 overflow-y-auto border-r pr-2">
					{#each data.documents as d (d.id)}
						<Item.Root
							class="hover:bg-muted cursor-pointer text-left {selected === d.id ? 'bg-muted' : ''}"
						>
							{#snippet child({ props })}
								<button type="button" onclick={() => select(d.id)} {...props}>
									<Item.Content>
										<Item.Title class="w-full truncate">{d.name}</Item.Title>
										<Item.Description class="flex items-center gap-1.5">
											{kb(d.size_bytes)}
											<Badge variant="outline">{humanizeEnum(d.status)}</Badge>
										</Item.Description>
									</Item.Content>
								</button>
							{/snippet}
						</Item.Root>
					{:else}
						<p class="text-muted-foreground p-2 text-sm">No documents yet.</p>
					{/each}
				</div>
				<div class="flex-1">
					{#if selectedDoc}
						<iframe
							title={selectedDoc.name}
							src="/documents/{selectedDoc.id}/content"
							class="h-full w-full rounded-md border"
						></iframe>
					{:else}
						<div
							class="text-muted-foreground flex h-full items-center justify-center rounded-md border border-dashed"
						>
							Select document to display
						</div>
					{/if}
				</div>
			</div>

			<Dialog.Root bind:open={uploadOpen}>
				<Dialog.Content>
					<Dialog.Header>
						<Dialog.Title>Upload document</Dialog.Title>
					</Dialog.Header>
					<form
						method="POST"
						action="?/upload"
						enctype="multipart/form-data"
						use:enhance={() =>
							async ({ result, update }) => {
								await update();
								if (result.type === 'success') uploadOpen = false;
							}}
						class="flex flex-col gap-4"
					>
						<Field>
							<FieldLabel for="file">File</FieldLabel>
							<input
								id="file"
								name="file"
								type="file"
								accept="application/pdf,.pdf"
								required
								class="text-sm file:mr-3 file:rounded-md file:border-0 file:bg-primary file:px-3 file:py-1 file:text-primary-foreground"
							/>
						</Field>
						{#if form?.error}<p class="text-destructive text-sm">{form.error}</p>{/if}
						<Dialog.Footer>
							<Button type="submit">Upload</Button>
						</Dialog.Footer>
					</form>
				</Dialog.Content>
			</Dialog.Root>

			<Dialog.Root bind:open={dialogOpen}>
				<Dialog.Content class="sm:max-w-4xl">
					<Dialog.Header>
						<Dialog.Title>
							{selectedDoc?.name}: {result?.status === 'FAILED' ? 'parsing error' : 'parsed output'}
						</Dialog.Title>
					</Dialog.Header>
					<CodeBlock
						label={result?.status === 'FAILED' ? 'Error' : 'Output'}
						content={result?.status === 'FAILED' ? result.error : result?.output}
					/>
				</Dialog.Content>
			</Dialog.Root>
		</Card.Content>
	</Card.Root>
</div>
