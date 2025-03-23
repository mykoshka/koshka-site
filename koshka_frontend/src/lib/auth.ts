export function getAuthHeaders() {
    const token = localStorage.getItem("auth_token");
    return token ? { Authorization: `Bearer ${token}` } : {};
}

export function getUserRole(): string {
    return localStorage.getItem("role") || "user"; // Default role is "user"
}

export function isAuthenticated(): boolean {
    return !!localStorage.getItem("auth_token");
}

export function isAdmin(): boolean {
    return getUserRole() === "admin";
}