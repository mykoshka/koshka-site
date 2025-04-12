import {type Writable, type Readable, writable, readonly} from 'svelte/store';
import type { ISections } from '$lib/types/menu'

export const menuOptions: Writable<ISections[]> = writable([]);

export const menuOptionsRO:Readable<ISections[]> = readonly(menuOptions);