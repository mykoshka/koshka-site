<script lang="ts">
    import type { PetProfile } from "../types/types";
    import { createEventDispatcher } from "svelte";

    export let pet: PetProfile;
    const dispatch = createEventDispatcher();

    function openPetProfile() {
        console.log("📌 Pet clicked:", pet);
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
        class="bg-white shadow-md rounded-lg p-4 cursor-pointer transition-transform transform hover:scale-105 focus:outline-none focus:ring focus:ring-primary"
        on:click={openPetProfile}
        on:keydown={(event) => event.key === "Enter" && openPetProfile()}
>
    <img src={"http://127.0.0.1:3000/images/" + pet.picture || "/default-pet.png"} alt={pet.name} class="w-full h-32 object-cover rounded-md" />
    <h3 class="text-lg font-semibold mt-2">{pet.name}</h3>
    <p class="text-sm text-gray-600">{pet.tag_id ? `Tag: ${pet.tag_id}` : "No Tag Assigned"}</p>
</div>
