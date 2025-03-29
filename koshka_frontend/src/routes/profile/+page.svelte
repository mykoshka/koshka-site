<script lang="ts">
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { getAuthHeaders } from "$lib/auth";
    import PetCard from "$lib/components/PetCard.svelte";
    import PetProfileModal from "$lib/components/PetProfileModal.svelte";
    import type { UserProfile, PetProfile } from "$lib/types/types";
    import { type GridOptions, type ColDef } from "ag-grid-community";
    import "ag-grid-community/styles/ag-grid.css";
    import "ag-grid-community/styles/ag-theme-alpine.css";

    let selectedPet: PetProfile | null = null;
    let profile: UserProfile | null = null;
    let error = '';

    // ✅ Define columns correctly
    const columnDefs: ColDef[] = [
        { field: "product_sku", headerName: "Product SKU", sortable: true, filter: true },
        { field: "registration_date", headerName: "Purchase Date", sortable: true, filter: true },
        { field: "tag_id", headerName: "Tag ID", sortable: true, filter: true }
    ];

    let gridOptions: GridOptions = {
        columnDefs,
        rowData: profile?.purchase_history,
        pagination: true,
        paginationPageSize: 5
    };

    async function fetchProfile() {
        const token = getAuthHeaders();
        if (!token) {
            error = 'Unauthorized. Redirecting to login...';
            setTimeout(() => goto('/login'), 2000);
            return;
        }

        try {
            const res = await fetch('http://localhost:3000/api/profile', {
                headers: token
            });
            if (res.ok) {
                const data = await res.json();
                profile = { ...data, pets: data.pets || [], purchase_history: data.purchase_history || [] }; // Ensure pets is always an array

            } else {
                error = 'Failed to fetch profile. Redirecting to login...';
                setTimeout(() => goto('/login'), 2000);
            }
        } catch (e) {
            console.error("Failed to load profile.");
            error = 'Network error. Please try again later.';
        }
    }

    function viewPetProfile(event) {
        if (!event.detail.pet) {
            console.error("❌ ERROR: `event.detail.pet` is undefined in Profile.svelte");
            return;
        }
        selectedPet = event.detail.pet;
    }

    function closeModal() {
        console.log("🔒 Closing Modal");
        selectedPet = null;
    }

    onMount(fetchProfile);
</script>

<div class="flex flex-col items-center h-fit p-6 bg-gray-100">
    <div class="w-full max-w-2xl p-6 bg-white rounded-lg shadow-md">
        {#if profile}
            <h2 class="text-2xl font-bold text-gray-700 pb-5">Profile</h2>
            <p><strong>Member Since:</strong> {profile.joined_on}</p>
            <p><strong>Name:</strong> {profile.name}</p>
            <p><strong>Email:</strong> {profile.email}</p>
            <p><strong>Mobile:</strong> {profile.mobile_number}</p>
            <p><strong>Address:</strong> {profile.address}</p>
            <button class="bg-green-50 flex shadow-md rounded-lg mt-5 p-4 cursor-pointer transition-transform transform hover:scale-105 focus:outline-none focus:ring focus:ring-primary border-green-500 hover:bg-green-500">Update Password</button>
        {/if}
        {#if error}
            <p class="text-red-500">{error}</p>
        {/if}
    </div>

    <div class="bg-gray-200 py-6 px-30">
        <h2 class="text-xl font-semibold mt-10">Linked Pets</h2>
        {#if profile?.pets?.length == 0}
            <p class="info text-gray-500 py-5">No pets linked to this account.</p>
        {:else if profile?.pets?.length > 2}
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-6 mt-4 w-full max-w-full">
                <PetCard pet={profile?.pets[0]} on:view={viewPetProfile} />
            </div>
        {:else}
            <div class="grid grid-cols-1 sm:grid-cols-1 gap-6 mt-4 w-full max-w-2xl">
                {#each profile?.pets as pet}
                    <PetCard {pet} on:view={viewPetProfile} />
                {/each}
            </div>
        {/if}
    </div>

    <!-- ✅ PURCHASE HISTORY TABLE -->
    <h2 class="text-xl font-semibold mt-14">Purchase History</h2>
    {#if profile?.purchase_history.length > 0}
        <div class="ag-theme-alpine w-full h-64 mt-4" bind:this={gridOptions.api}></div>
    {:else}
        <p class="text-gray-500 py-5">No purchase history found.</p>
    {/if}
</div>
{#if selectedPet}
    <PetProfileModal selectedPet={selectedPet} on:close={closeModal} />
{/if}