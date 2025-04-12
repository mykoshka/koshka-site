<script lang="ts">
    import {onMount} from 'svelte';
    import {goto} from '$app/navigation';
    import {getAuthHeaders, storeTokens} from "$lib/auth";
    // import SideMenu from "$lib/components/profile/SideMenu.svelte";
    import PetCard from "$lib/components/profile/PetCard.svelte";
    import PetProfileModal from "$lib/components/profile/PetProfileModal.svelte";
    //    import type { NonNullableFields, UserProfile, PetProfile } from "$lib/types/profile";
    import type {PetProfile, UserProfile} from "$lib/types/profile";
    import {EmptyUserProfile} from "$lib/types/profile";
    import type {ISections} from '$lib/types/menu'
    import {menuOptions} from '$lib/stores/menuStore'
    import {type ColDef, type GridOptions} from "ag-grid-community";
    import "ag-grid-community/styles/ag-grid.css";
    import "ag-grid-community/styles/ag-theme-alpine.css";

    /* ------------------------------------------------
    Variable Section
    ------------------------------------------------ */
    let selectedPet: PetProfile | null = null;
    let profile: UserProfile = EmptyUserProfile;
    let gridOptions: GridOptions;
    let error = '';

    /* ------------------------------------------------
    AG Grid Section
    ------------------------------------------------ */
    if (profile.purchase_history) {
        // ✅ Define columns correctly
        const columnDefs: ColDef[] = [
            {field: "product_sku", headerName: "Product SKU", sortable: true, filter: true},
            {field: "registration_date", headerName: "Purchase Date", sortable: true, filter: true},
            {field: "tag_id", headerName: "Tag ID", sortable: true, filter: true}
        ];

        gridOptions = {
            columnDefs,
            rowData: profile?.purchase_history,
            pagination: true,
            paginationPageSize: 5
        };
    } else {
        gridOptions = {};
    }
    /* ------------------------------------------------
    Data fetch Section
    ------------------------------------------------ */

    async function fetchProfile() {
        const token = getAuthHeaders();
        const MenuItems: Array<ISections> = [];
        if (!token) {
            error = 'Unauthorized. Redirecting to login...';
            setTimeout(() => goto('/login'), 2000);
            return;
        }
        MenuItems.push({id: "profile", label: "Profile"})
        try {
            const res = await fetch('http://localhost:3000/api/profile', {
                headers: token
            });
            if (res.ok) {
                const body = await res.json()
                const data = body["profile"];
                const newToken = body["new_tokens"];
                if ( newToken != null ) {
                    storeTokens(newToken.new_token, newToken.new_refresh)
                }
                profile = {...data, pets: data.pets || [], purchase_history: data.purchase_history || []}; // Ensure pets is always an array
                if (data.pets != null) {
                    MenuItems.push({id: "pets", label: "Pets"})
                }
                if (data.purchase_history != null) {
                    MenuItems.push({id: "purchases", label: "Purchases"})
                }
            } else {
                error = 'Failed to fetch profile. Redirecting to login...';
                setTimeout(() => goto('/login'), 2000);
            }
           menuOptions.set(MenuItems)
        } catch (e) {
            console.error("Failed to load profile.", e);
            error = 'Network error. Please try again later.';
        }
    }

    /* ------------------------------------------------
    Detailed Pet Profile Section
    ------------------------------------------------ */

    function viewPetProfile(event: CustomEvent) {
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

    /* ------------------------------------------------
    Final steps Section
    ------------------------------------------------ */

    onMount(async () => {
        await fetchProfile();
    });

</script>

<!--<SideMenu></SideMenu>-->
{#if profile}
    <div class="flex flex-col items-center h-fit w-full p-6 bg-white">
        <div id="profile" class="w-4/5 my-3.5 p-6 bg-gray-50 rounded-lg shadow-md">
            {#if profile}
                <h2 class="text-2xl font-bold text-gray-700 pb-5">Profile</h2>
                <p><strong>Member Since:</strong> {profile.joined_on}</p>
                <p><strong>Name:</strong> {profile.name}</p>
                <p><strong>Email:</strong> {profile.email}</p>
                <p><strong>Mobile:</strong> {profile.mobile_number}</p>
                <p><strong>Address:</strong> {profile.address}</p>
                <button class="bg-green-50 flex shadow-md rounded-lg mt-5 p-4 cursor-pointer transition-transform transform hover:scale-105 focus:outline-none focus:ring focus:ring-primary border-green-500 hover:bg-green-500">
                    Update Password
                </button>
            {/if}
            {#if error}
                <p class="text-red-500">{error}</p>
            {/if}
        </div>

        <div id="pets" class="w-4/5 my-3.5 p-6 bg-gray-50 rounded-lg shadow-md">
            <h2 class="text-2xl font-bold text-gray-700 pb-5">Linked Pets</h2>
            {#if profile.pets == null}
                <p class="info text-gray-500 p-6">No pets linked to this account.</p>
            {:else if profile.pets.length > 2}
                <div class="grid grid-cols-1 sm:grid-cols-1 gap-6 mt-4 w-full max-w-full">
                    <PetCard pet={profile.pets[0]} col=1 on:view={viewPetProfile}/>
                </div>
            {:else}
                <div class="grid grid-cols-1 content-center sm:grid-cols-2 gap-6 mt-4 w-full">
                    {#each profile?.pets as pet(pet.id)}
                        <PetCard {pet} col=2 on:view={viewPetProfile}/>
                    {/each}
                </div>
            {/if}
        </div>

        <!-- ✅ PURCHASE HISTORY TABLE -->
        <div id="pets" class="w-4/5 my-3.5 p-6 bg-gray-50 rounded-lg shadow-md">
            <h2 id="purchases" class="text-2xl font-bold text-gray-700 pb-5">Purchase History</h2>
            {#if profile?.purchase_history}
                <p class="text-gray-500 py-5">No purchase history found.</p>
            {:else}
                <div class="ag-theme-alpine w-full h-64 mt-4" bind:this={gridOptions.api}></div>
            {/if}
        </div>
    </div>
{/if}
{#if selectedPet}
    <PetProfileModal selectedPet={selectedPet} on:close={closeModal}/>
{/if}