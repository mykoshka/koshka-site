import {type Writable, writable} from 'svelte/store';

export interface IuserProfile {
    isLoggedIn: boolean;
    role: string | null;
}
export const userProfile: Writable<IuserProfile> = writable({ isLoggedIn: false, role: null });

export function updateUserProfile(newProfile: IuserProfile): void {
    userProfile.update((current) => ({ ...current, ...newProfile }));
}