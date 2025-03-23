<script lang="ts">
    import { onMount } from "svelte";
    import axios from "axios";
    import { getAuthHeaders, isAdmin, isAuthenticated } from "$lib/auth";
    import { goto } from "$app/navigation";

    let name = "";
    let permissions = "read-only"; // Default permission
    let apiKeys: { id: string; name: string; permissions: string[] }[] = [];
    let errorMessage = "";

    async function fetchApiKeys() {
        try {
            const res = await axios.get("/api/admin/api-keys", { headers: getAuthHeaders() });
            apiKeys = res.data; // Assuming API returns an array of { id, name, permissions }
        } catch (error) {
            errorMessage = "Failed to fetch API keys.";
        }
    }

    async function createApiKey() {
        try {
            const res = await axios.post("/api/admin/api-keys", { name, permissions: [permissions] }, { headers: getAuthHeaders() });
            apiKeys.push({ id: res.data.api_key, name, permissions: [permissions] }); // Add to UI
            name = ""; // Reset input
        } catch (error) {
            errorMessage = "Failed to create API key.";
        }
    }

    async function deleteApiKey(id: string) {
        try {
            await axios.delete(`/api/admin/api-keys/${id}`, { headers: getAuthHeaders() });
            apiKeys = apiKeys.filter((key) => key.id !== id); // Remove from UI
        } catch (error) {
            errorMessage = "Failed to delete API key.";
        }
    }

    onMount(() => {
        if (!isAuthenticated() || !isAdmin()) goto("/login");
        fetchApiKeys();
    });
</script>

<div class="p-6">
    <h1 class="text-2xl font-semibold">Admin: API Key Management</h1>

    <div class="mt-4">
        <input bind:value={name} placeholder="API Key Name" class="border p-2 rounded w-full" />
        <select bind:value={permissions} class="border p-2 rounded w-full mt-2">
            <option value="read-only">Read-Only</option>
            <option value="write">Write</option>
            <option value="admin">Admin</option>
        </select>
        <button class="mt-4 px-4 py-2 bg-blue-500 text-white rounded-lg" on:click={createApiKey}>Create API Key</button>
    </div>

    <h2 class="text-xl mt-6">Existing API Keys</h2>
    {#if apiKeys.length > 0}
        <table class="w-full mt-4 border">
            <tr class="bg-gray-100">
                <th class="p-2">ID</th>
                <th class="p-2">Name</th>
                <th class="p-2">Permissions</th>
                <th class="p-2">Actions</th>
            </tr>
            {#each apiKeys as key}
                <tr class="border-t">
                    <td class="p-2 font-mono">{key.id}</td>
                    <td class="p-2">{key.name}</td>
                    <td class="p-2">{key.permissions.join(", ")}</td>
                    <td class="p-2">
                        <button class="text-red-500" on:click={() => deleteApiKey(key.id)}>Delete</button>
                    </td>
                </tr>
            {/each}
        </table>
    {:else}
        <p class="text-gray-500 mt-4">No API keys found.</p>
    {/if}

    {#if errorMessage}
        <p class="text-red-500 mt-4">{errorMessage}</p>
    {/if}
</div>