<script lang="ts">
	import type { User, Issue, Document } from '$lib/api';
	import { getDocumentResult, reprocessDocument, deleteDocument } from '$lib/api';
	import { clientApi } from '$lib/api-client';
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { ButtonGroup } from '$lib/components/ui/button-group/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Item from '$lib/components/ui/item/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import ApprovalPanel from '$lib/components/approval/approval-panel.svelte';
	import CodeBlock from '$lib/components/code-block.svelte';
	import { kb, humanizeEnum } from '$lib/fmt';
	import { FileText, Maximize, Minimize, RotateCcw, TriangleAlert, Upload } from '@lucide/svelte';
	import { Field, FieldLabel } from '$lib/components/ui/field';
	import { Trash2 } from '@lucide/svelte';
	import { enhance } from '$app/forms';
	import { invalidateAll } from '$app/navigation';
	import { toast } from '$lib/toast';

	type Props = {
		selectedDocId: string | null;
		onSelect: (id: string) => void;
		formError: string | undefined;
		data: {
			user: User | null;
			issue: Issue;
			documents: Document[];
		};
	};
	let { data, formError, selectedDocId, onSelect }: Props = $props();
	let fullScreen = $state(false);

	const selectedDoc = $derived(data.documents.find((d) => d.id === selectedDocId) ?? null);
	type DocParserResult = {
		status: 'QUEUED' | 'PROCESSING' | 'DONE' | 'FAILED';
		error?: string | null;
		output?: string | null;
	};
	let result = $state<DocParserResult | null>(null);
	let dialogOpen = $state(false);
	let uploadOpen = $state(false);

	let restartTrigger = $state(0);

	$effect(() => {
		const id = selectedDocId;
		void restartTrigger;
		result = null;
		if (!id) return;
		let stopped = false;
		let timer: ReturnType<typeof setTimeout>;
		const poll = async () => {
			const { data, error } = await getDocumentResult({
				client: clientApi,
				path: { documentId: id }
			});
			if (stopped || error || !data) return;
			result = data;
			if (result.status === 'QUEUED' || result.status === 'PROCESSING') {
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
			QUEUED: 'Queued for parsing...',
			PROCESSING: 'Parsing...',
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

	async function restartParser() {
		const id = selectedDocId;
		if (!id) return;
		await reprocessDocument({ client: clientApi, path: { documentId: id } });
		++restartTrigger;
	}

	let deleteTarget = $state<Document | null>(null);
	let deleteOpen = $state(false);

	function askDelete(d: Document) {
		deleteTarget = d;
		deleteOpen = true;
	}

	async function confirmDelete() {
		const t = deleteTarget;
		if (t) {
			const { error } = await deleteDocument({ client: clientApi, path: { documentId: t.id } });
			if (error) {
				toast.error('Failed to delete document');
			} else {
				toast.success('Document deleted');
				await invalidateAll();
			}
		}
		deleteOpen = false;
		deleteTarget = null;
	}
</script>

<svelte:window
	onkeydown={(e) => {
		if (e.key === 'Escape' && !dialogOpen && !uploadOpen) fullScreen = false;
	}}
/>

<div
	class={fullScreen
		? 'bg-background fixed inset-0 z-40 flex flex-col gap-3 p-6'
		: 'flex flex-col gap-3'}
>
	<Tabs.Root value="preview" class="contents">
		<div class="flex items-center justify-between gap-2 text-sm">
			<div class="flex items-center gap-2">
				{#if fullScreen}
					<Tooltip.Provider>
						<Tooltip.Root>
							<Tooltip.Trigger>
								{#snippet child({ props })}
									<h1 {...props} class="text-muted-foreground border-b-2 cursor-default">
										{data.issue.title}
									</h1>
								{/snippet}
							</Tooltip.Trigger>
							<Tooltip.Content class="text-xl"
								>{data.issue.description ?? 'No description'}</Tooltip.Content
							>
						</Tooltip.Root>
					</Tooltip.Provider>
					|
				{/if}
				<div class="flex items-center gap-2 min-w-sm">
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
								<Button size="sm" variant="outline" onclick={restartParser}>
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
			</div>
			{#if selectedDoc}
				<Tabs.List>
					<Tabs.Trigger value="preview">Preview</Tabs.Trigger>
					<Tabs.Trigger value="approval">Approval</Tabs.Trigger>
				</Tabs.List>
			{/if}
			<div class="flex items-center gap-2">
				<Button
					size="sm"
					variant="outline"
					onclick={() => (fullScreen = !fullScreen)}
					aria-label="Toggle fullscreen"
				>
					{#if fullScreen}<Minimize />{:else}<Maximize />{/if}
				</Button>
				<Button size="sm" onclick={() => (uploadOpen = true)}>
					<Upload />
					Upload
				</Button>
			</div>
		</div>
		<div class="flex gap-4 {fullScreen ? 'min-h-0 flex-1' : 'h-[70vh]'}">
			<div class="w-64 shrink-0 space-y-1 overflow-y-auto border-r pr-2">
				{#each data.documents as d (d.id)}
					<Item.Root class="hover:bg-muted {selectedDocId === d.id ? 'bg-muted' : ''}">
						<button
							type="button"
							onclick={() => onSelect(d.id)}
							class="flex min-w-0 flex-1 cursor-pointer flex-col items-start gap-0.5 text-left"
						>
							<Item.Title class="w-full truncate">{d.name}</Item.Title>
							<Item.Description class="flex items-center gap-1.5">
								{kb(d.size_bytes)}
								<Badge variant="outline">{humanizeEnum(d.status)}</Badge>
							</Item.Description>
						</button>
						<Item.Actions>
							<Button
								size="icon-sm"
								variant="ghost"
								aria-label="Delete document"
								onclick={() => askDelete(d)}
							>
								<Trash2 />
							</Button>
						</Item.Actions>
					</Item.Root>
				{:else}
					<p class="text-muted-foreground p-2 text-sm">No documents yet.</p>
				{/each}
			</div>
			<div class="flex flex-1 flex-col">
				{#if selectedDoc}
					<Tabs.Content value="preview" class="mt-0 min-h-0 flex-1">
						<iframe
							title={selectedDoc.name}
							src="/documents/{selectedDoc.id}/content"
							class="h-full w-full rounded-md border"
						></iframe>
					</Tabs.Content>
					<Tabs.Content value="approval" class="mt-0 min-h-0 flex-1 overflow-hidden">
						<ApprovalPanel document={selectedDoc} currentUser={data.user} />
					</Tabs.Content>
				{:else}
					<div
						class="text-muted-foreground flex h-full items-center justify-center rounded-md border border-dashed"
					>
						Select document to display
					</div>
				{/if}
			</div>
		</div>
	</Tabs.Root>
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
			{#if formError}<p class="text-destructive text-sm">{formError}</p>{/if}
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

<AlertDialog.Root bind:open={deleteOpen}>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Delete document?</AlertDialog.Title>
			<AlertDialog.Description>
				“{deleteTarget?.name}” and its parsing result and any approval processes will be permanently
				deleted. This can’t be undone.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				class={buttonVariants({ variant: 'destructive' })}
				onclick={confirmDelete}
			>
				Delete
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
