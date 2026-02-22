type ToastType = 'success' | 'error' | 'info';

interface Toast {
  id: number;
  type: ToastType;
  message: string;
}

let toasts = $state<Toast[]>([]);
let nextId = 0;

function addToast(type: ToastType, message: string, duration = 4000) {
  const id = nextId++;
  toasts = [...toasts, { id, type, message }];
  if (toasts.length > 3) toasts = toasts.slice(-3);
  setTimeout(() => removeToast(id), duration);
}

function removeToast(id: number) {
  toasts = toasts.filter((t) => t.id !== id);
}

export const toast = {
  success: (msg: string) => addToast('success', msg),
  error: (msg: string) => addToast('error', msg),
  info: (msg: string) => addToast('info', msg),
};

export function getToasts(): Toast[] {
  return toasts;
}

export { removeToast };
