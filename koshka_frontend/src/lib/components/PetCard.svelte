<script lang="ts">
    import type { PetProfile } from "$lib/types/types";
    import { createEventDispatcher } from "svelte";

    export let pet: PetProfile;
    const dispatch = createEventDispatcher();

    function openPetProfile() {
        if (!pet) {
            console.error("❌ ERROR: `pet` is undefined in PetCard.svelte");
            return;
        }
        dispatch("view", { pet });
    }
</script>

<div
        role="button"
        tabindex="0"
        class="bg-white flex shadow-md rounded-lg p-4 cursor-pointer transition-transform transform hover:scale-105 focus:outline-none focus:ring focus:ring-primary"
        on:click={openPetProfile}
        on:keydown={(event) => event.key === "Enter" && openPetProfile()}
>
    <div class="shrink mr-5">
        <img src={"http://127.0.0.1:3000/images/" + pet.picture || "/default-pet.png"} alt={pet.name} class="h-32 w-32 object-cover rounded-md" />
    </div>
    <div class="col-span-2">
        <div class="text-sm text-gray-300 flex-grow w-full">Koshka ID: {pet.id}</div>
        <h3 class="text-lg font-semibold mb-1">{pet.name}</h3>
            <div class="grid grid-row-2 grid-cols-2 grid-flow-row-dense gap-x-4">
                <p class="text-sm text-gray-600">Date of Birth: {pet.date_of_birth != "1900-01-01" ? `${pet.date_of_birth}` : "No DOB Entered" }</p>
                <p class="text-sm text-gray-600">City Licence #: {pet.city_license ? `${pet.city_license}` : "No City Licence"}</p>
                <p class="text-sm text-gray-600">Neutered: {pet.neutered}</p>
                <p class="text-sm text-gray-600">Vaccinated: {pet.vaccinated}</p>
            </div>
        <p class="text-sm text-gray-600 mt-3">{pet.tag_id ? `Tag: ${pet.tag_id}` : "No Tag Assigned"}</p>
    </div>
</div>
