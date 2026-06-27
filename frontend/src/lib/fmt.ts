export function humanizeEnum(value: string): string {
  return value
    .split('_')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ');
}

export function kb(bytes: number) {
	return `${(bytes / 1024).toFixed(1)} KB`;
}