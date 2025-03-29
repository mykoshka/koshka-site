<script>
    import icons from "./icons.js";
    import {createEventDispatcher} from "svelte";
    //look at all these beautiful options
    // Buttons text, set any to "" to remove that button
    export let descriptionText = "Drag or click to upload";
    // The file upload input element
    export let input = null;
    //Fixed uploader position?
    export let fixed = true;
    //Files from the file input and the drag zone
    export let inputFiles = [];
    export let dragZoneFiles = [];
    //Trigger file upload
    export let trigger = () => input.click();
    //External method to get the current files at any time
    export const getFiles = () => files;
    // Called when maxuploads is reached or the done button is clicked
    export let callback = () => {
    };
    //Called when the "Done" button is clicked
    export let doneCallback = () => {
    };
    // Drag zone element
    export let dragZone = null;
    //Maximum files that can be uploaded
    export let maxFiles = 3;
    // When the maximum files are uploaded
    export let maxFilesCallback = () => {
    };
    //Show a list of files + icons?
    export let listFiles = true;

    $: files = [...inputFiles, ...dragZoneFiles];
    $: {
        if (files.length >= maxFiles) {
            maxFilesCallback(files, maxFiles);
            dispatch("change", files)
            callback(files);
        }
    }

    const dispatch = createEventDispatcher();

    function dragover(e) {
        e.preventDefault();
        dispatch("dragover", e);
    }

    function dragenter(e) {
        e.preventDefault();
        dragZone.classList.add("dragging");
        dispatch("dragenter", e);
    }

    function dragleave(e) {
        e.preventDefault();
        dragZone.classList.remove("dragging");
        dispatch("dragleave", e);
    }

    function drop(e) {
        e.preventDefault();
        dragZone.classList.remove("dragging");
        dragZoneFiles.push(...[...e.dataTransfer.items].filter(i => i.kind === "file").map(i => i.getAsFile()))
        dragZoneFiles = [...dragZoneFiles];
        dispatch("drop", e);
        dispatch("change", files)
    }

    function formatBytes(a, b = 2, k = 1024) {
        let d = Math.floor(Math.log(a) / Math.log(k));
        return 0 == a
            ? "0 Bytes"
            : parseFloat((a / Math.pow(k, d)).toFixed(Math.max(0, b))) +
            " " +
            ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"][d];
    }

    function getIcon(filename) {
        if (!filename) {
            return icons.default
        }
        return Object.entries(icons).find(
            i => i[0]
                .split(",")
                .includes(filename.split(".").slice(-1)[0])
        )?.[1] || icons.default;
    }

    function inputChanged() {
        inputFiles = [...inputFiles, ...input.files]
    }

    function del(file) {
        if (idx(file, inputFiles) !== null) {
            inputFiles.splice(idx(file, inputFiles), 1);
            inputFiles = [...inputFiles];
            return;
        }
        if (idx(file, dragZoneFiles) !== null) {
            dragZoneFiles.splice(idx(file, dragZoneFiles), 1);
            dragZoneFiles = [...dragZoneFiles];
            return;
        }
        return console.log(idx(file, inputFiles), idx(file, dragZoneFiles))

        function idx(item, arr) {
            let i = arr.findIndex(i => i === item);
            if (i < 0) {
                return null
            } else {
                return i
            }
        }
    }

    function openFile(file) {
        window.open(URL.createObjectURL(file), "filewin");
    }
</script>
<div bind:this={dragZone} on:dragover={dragover} on:drop={drop} on:dragenter={dragenter} on:dragleave={dragleave}
     class={`${fixed ? 'dropfixed' : ''} fileUploader dragzone`}>
    {#if files.length !== maxFiles}
        {#if listFiles}
            <ul>
                {#each files.slice(0, maxFiles) as file}
                    <li on:click={() => openFile(file)}>
                <span class="icon">
                    <span class="fileicon">{@html getIcon(file.name)}</span>
                    <span class="deleteicon" on:click|stopPropagation={() => del(file)}><svg
                            xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
                            aria-hidden="true" role="img" class="iconify iconify--ph" width="32" height="32"
                            preserveAspectRatio="xMidYMid meet" viewBox="0 0 256 256"><path fill="currentColor"
                                                                                            d="M216 48h-40v-8a24.1 24.1 0 0 0-24-24h-48a24.1 24.1 0 0 0-24 24v8H40a8 8 0 0 0 0 16h8v144a16 16 0 0 0 16 16h128a16 16 0 0 0 16-16V64h8a8 8 0 0 0 0-16ZM96 40a8 8 0 0 1 8-8h48a8 8 0 0 1 8 8v8H96Zm96 168H64V64h128Zm-80-104v64a8 8 0 0 1-16 0v-64a8 8 0 0 1 16 0Zm48 0v64a8 8 0 0 1-16 0v-64a8 8 0 0 1 16 0Z"></path></svg></span>
                </span>
                        <span class="filename">{file.name}</span>
                        <span class="filesize">{formatBytes(file.size)}</span>
                    </li>
                {/each}
            </ul>
        {/if}
    {/if}
</div>
<input type="file" hidden bind:this={input} on:change={inputChanged} multiple={maxFiles > 1}>

<style>
    .dragzone {
        padding: 20px;
        border-radius: 6px;
        border: 2px dashed #0001;
        display: flex;
        justify-content: center;
        align-items: center;
        flex-direction: column;
        transition: background-color .3s ease;
    }

    .dragzone ul {
        width: 60%;
        display: flex;
        flex-direction: column;
    }

    .dragzone li {
        transition: background-color .2s ease-in-out;
        list-style: none;
        display: flex;
        align-items: center;
        cursor: pointer;
        padding: 3px 8px;
        border-radius: 3px;
    }

    .dragzone li:hover {
        background: #0001;
    }

    .dragzone .filesize {
        flex: 1;
        text-align: right;
        opacity: .2;
        font-style: italic;
    }

    .dragzone li .filename {
        white-space: nowrap;
        overflow: hidden;
        width: 15ch;
        text-overflow: ellipsis;
        display: block;
        font-weight: 200;
        opacity: .8;
    }

    .dragzone.dragging {
        background: #0662;
    }

    .dragzone svg:not(li svg) {
        width: 15vw;
        height: 15vw;
        color: #777;
        opacity: .6;
    }

    .dragzone li .icon {
        width: 20px;
        margin-right: 10px;
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .dragzone li .icon :global(svg) {
        width: 20px;
        color: #333;
    }

    .deleteicon {
        display: none;
    }

    .dragzone li:hover .fileicon {
        display: none;
    }

    .dragzone li:hover .deleteicon {
        display: block;
        cursor: pointer;
    }

    .dragzone .doneText {
        font-size: 1.3rem;
        color: #333;
        opacity: .5;
        font-weight: 200;
        font-style: italic;
        margin-top: 2rem;
    }

    .dragzone.dropfixed {
        position: absolute;
        width: 100%;
        height: 100%;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
    }

    .dragzone .text {
        font-style: italic;
        color: #333;
        opacity: .5;
        font-weight: 200;
        padding-bottom: 3rem;
    }

    .buttons {
        width: 40%;
        display: flex;
    }

    .buttons button {
        margin: 0 5px;
        transition: background-color .2s ease;
        padding: .5rem 1rem;
        margin-bottom: 1rem;
        flex: 1;
        border: 1px solid #0001;
        background: transparent;
        cursor: pointer;
    }

    .buttons button:hover {
        background: #0001;
    }
</style>