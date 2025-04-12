import {updateUserProfile} from '$lib/stores/authStore';

export function getAuthHeaders() {
    const token = localStorage.getItem("auth_token");
    const refresh = localStorage.getItem("refresh_token");

    const newToken = localStorage.getItem("new_auth");
    const newRefresh = localStorage.getItem("new_refresh");

    if ( newToken && newRefresh ) {storeTokens(newToken, newRefresh)}

    return token && refresh ? { Authorization: `Bearer ${token}`, Refresh: refresh } : {};
}

export function getUserRole(): string {
    return localStorage.getItem("role") || "user"; // Default role is "user"
}

export function storeTokens(bearer: string, refresh: string): void {
    localStorage.setItem("auth_token", bearer);
    localStorage.setItem("refresh_token", refresh);
    updateUserProfile({ isLoggedIn: true, role: null })
}

export function isAuthenticated(): boolean {
    const auth = !!localStorage.getItem("auth_token") && localStorage.getItem("auth_token") !== undefined
    const refresh = !!localStorage.getItem("refresh_token") && localStorage.getItem("refresh_token") !== undefined
    return auth && refresh;
}

export function isAdmin(): boolean {
    return getUserRole() === "admin";
}