<script lang="ts">
	import { page } from '$app/state';

	const { error, status } = page;
</script>

<div class="space-y-1">
	<h4 class="text-sm leading-none font-medium">{status}</h4>
	<p class="text-muted-foreground text-sm">
		{error?.message ?? 'Unknown error'}
	</p>
</div>
