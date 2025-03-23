import { writable } from 'svelte/store';

export interface IuserProfile {
    isLoggedIn: boolean;
    role: string | null;
}
export const userProfile = writable({ isLoggedIn: false, role: null });
