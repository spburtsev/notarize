import { toast as sonner, type ExternalToast } from 'svelte-sonner';
import Toast, { type ToastVariant } from '$lib/components/toast.svelte';

function show(variant: ToastVariant, message: string, description?: string, opts?: ExternalToast) {
	return sonner.custom(Toast, {
		...opts,
		componentProps: { variant, message, description }
	});
}

export const toast = {
	success: (message: string, description?: string, opts?: ExternalToast) =>
		show('success', message, description, opts),
	error: (message: string, description?: string, opts?: ExternalToast) =>
		show('error', message, description, opts),
	info: (message: string, description?: string, opts?: ExternalToast) =>
		show('info', message, description, opts),
	warning: (message: string, description?: string, opts?: ExternalToast) =>
		show('warning', message, description, opts),
	dismiss: sonner.dismiss
};
