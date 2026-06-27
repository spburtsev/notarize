export function humanizeEnum(value: string): string {
	return value
		.split('_')
		.map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
		.join(' ');
}

export function kb(bytes: number) {
	return `${(bytes / 1024).toFixed(1)} KB`;
}

type BadgeVariant = 'default' | 'secondary' | 'destructive' | 'outline';

export function processStatusVariant(status: string): BadgeVariant {
	switch (status) {
		case 'APPROVED':
			return 'default';
		case 'REJECTED':
			return 'destructive';
		case 'CANCELLED':
			return 'outline';
		default:
			return 'secondary'; // PENDING, IN_PROGRESS
	}
}

export function stepStatusVariant(status: string): BadgeVariant {
	switch (status) {
		case 'APPROVED':
			return 'default';
		case 'REJECTED':
			return 'destructive';
		case 'ACTIVE':
			return 'secondary';
		default:
			return 'outline'; // PENDING, SKIPPED
	}
}

export function shortId(id: string) {
	return id.slice(0, 8);
}

export function formatDateTime(iso: string) {
	return new Date(iso).toLocaleString();
}
