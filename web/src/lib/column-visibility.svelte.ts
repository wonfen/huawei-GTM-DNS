import { t } from '$lib/i18n.svelte';

const STORAGE_KEY = 'gtm-col-vis';

const DEFAULTS = {
  probes: true,
  name: true,
  records: true,
  type: false,
  ttl: false,
  line: true,
  status: true,
  health: true,
  dns_weight: true,
  created_at: false,
  updated_at: false,
  description: false,
};

export type ColVisKeys = keyof typeof DEFAULTS;

const COL_KEYS: ColVisKeys[] = [
  'probes', 'name', 'records', 'type', 'ttl', 'line',
  'status', 'health', 'dns_weight', 'created_at', 'updated_at', 'description',
];

export function getColLabels(): Record<ColVisKeys, string> {
  return Object.fromEntries(
    COL_KEYS.map((k) => [k, t(`col.${k}`)])
  ) as Record<ColVisKeys, string>;
}

// Keep backward compat — a static export that re-evaluates via the function
export const COL_LABELS = new Proxy({} as Record<ColVisKeys, string>, {
  get(_target, prop: string) {
    return t(`col.${prop}`);
  },
  ownKeys() {
    return COL_KEYS;
  },
  getOwnPropertyDescriptor() {
    return { configurable: true, enumerable: true };
  },
});

function load(): typeof DEFAULTS {
  if (typeof localStorage === 'undefined') return { ...DEFAULTS };
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored) return { ...DEFAULTS, ...JSON.parse(stored) };
  } catch {}
  return { ...DEFAULTS };
}

export const colVis = $state(load());

$effect.root(() => {
  $effect(() => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify({ ...colVis }));
  });
});
