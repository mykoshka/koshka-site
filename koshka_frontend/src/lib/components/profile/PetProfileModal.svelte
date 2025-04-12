<script lang="ts">
    import {createEventDispatcher, onMount} from 'svelte';
    import Switch from '$lib/components/design/Switch.svelte'
    import FileUpload from "$lib/components/design/FileUploaderNoButton.svelte";

    let newImage: unknown = null;
    const dispatch = createEventDispatcher();

    export let selectedPet = {};

    onMount(() => {
        pet = {...selectedPet};
    });

    function closeModal() {
        dispatch('close');
    }

    function handleImageUpload(event) {
        const file = event.target.files[0];
        if (file) {
            newImage = file;
        }
    }

    export let pet = {
        id: selectedPet.id,
        picture: selectedPet.picture,
        name: selectedPet.name,
        date_of_birth: selectedPet.date_of_birth,
        tag_id: selectedPet.tag_id,
        city_licence: selectedPet.city_licence,
        neutered: selectedPet.neutered,
        vaccinated: selectedPet.vaccinated
    };

    async function saveChanges() {
        const formData = new FormData();
        if (newImage) {
            formData.append('image', newImage);
        }
        Object.keys(pet).forEach(key => {
            if (key !== 'id') {
                formData.append(key, pet[key]);
            }
        });

        closeModal();
    }

</script>

{#if selectedPet}
    <div class="fixed inset-0 flex items-center h-screen z-50 justify-center Opacity50">
        <div class="bg-white p-6 rounded-lg w-full max-w-3/4">
            <h2 class="text-2xl font-bold mb-4">Edit {pet.name}'s Profile</h2>
            <div class="text-xs text-gray-400">Koshka ID: {pet.id}</div>
            <div id="dragdropzone" class="flex justify-center" on:drop={handleImageUpload}>
                <div class="relative w-full h-full min-h-50 py-10">
                    <!--
                        <label for="pet-image" class="block mb-2">Upload Image:</label>
                        <input class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full"
                        id="pet-image_old" type="file" accept="image/*" on:change={handleImageUpload} />
                    -->
                    <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2">
                        <div class="grid place-items-center">
                            <img id="imgUpload" src={"http://127.0.0.1:3000/images/" + pet.picture} alt="Pet Image"
                             class="rounded-4xl w-40 h-40 object-cover"/>
                        </div>
                        <div class="dragzoneText">Drag and Drop to upload image</div>
                    </div>
                    <FileUpload fixed="true" listFiles="false" buttonText="upload" descriptionText="Drag and drop to upload" doneButtonText="" dragZone="dragdropzone" maxFiles="1" class="" on:change={() => console.log("Files changed")} callback={() => alert("Thanks")}>
                    </FileUpload>


                </div>
            </div>
            <div class="grid grid-cols-2">
                <div class="px-2">
                    <label for="pet-name" class="block mt-4">Name:</label>
                    <input id="pet-name" type="text" bind:value={pet.name} class="rounded-md border p-2 w-full"/>
                </div>
                <div class="px-2">
                    <label for="date_of_birth" class="block mt-4">Date of Birth:</label>
                    <input id="date_of_birth" type="date" bind:value={pet.date_of_birth}
                           class="rounded-md border p-2 w-full"/>
                </div>
                <div class="px-2">
                    <label for="tag_id" class="block mt-4">Tag ID:</label>
                    <input id="tag_id" type="text" bind:value={pet.tag_id} class="rounded-md border p-2 w-full"/>
                </div>
                <div class="px-2">
                    <label for="city_licence" class="block mt-4">City License:</label>
                    <input id="city_licence" type="text" bind:value={pet.city_licence}
                           class="rounded-md border p-2 w-full"/>
                </div>
                <div class="p-3">
                    <Switch bind:value={pet.neutered} label="Neutered" fontSize={16} design="slider"/>
                </div>
                <div class="p-3">
                    <Switch bind:value={pet.vaccinated} label="Vaccinated" fontSize={16} design="slider"/>
                </div>
            </div>

            <div class="flex justify-between mt-4">
                <button class="bg-gray-400 px-4 py-2 rounded hover:bg-red-600 hover:text-gray-100"
                        on:click={closeModal}>Cancel
                </button>
                <button class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-green-600 hover:text-gray-100"
                        on:click={saveChanges}>Save
                </button>
            </div>
        </div>
    </div>
{/if}

<style>
    label {
        padding: 0.75rem 1.5rem;
        border-radius: 0.5rem;
        cursor: pointer;
    }
    .dragzoneText {
        font-style: italic;
        color: #333;
        opacity: .5;
        font-weight: 200;
    }
    .Opacity50 {
        background-color:rgba(0, 0, 0, 0.8);
    }
</style>