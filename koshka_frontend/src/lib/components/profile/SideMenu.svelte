<script lang="ts">
    import {onMount} from 'svelte';
    import { menuOptionsRO } from '$lib/stores/menuStore';

    let activeSection = '';

    function handleScroll() {
        for (const section of $menuOptionsRO) {
            const el = document.getElementById(section.id);
            if (el && el.getBoundingClientRect().top <= 100) {
                activeSection = section.id;
            }
        }
    }

    onMount(() => {
        console.log($menuOptionsRO)
        activeSection = "profile"
        window.addEventListener('scroll', handleScroll);
        return () => window.removeEventListener('scroll', handleScroll);
    });
</script>

<style>
    .menu-item {
        padding: 10px;
        display: block;
        text-decoration: none;
        color: black;
        border-left: 3px solid transparent;
        transition: all 0.3s;
    }

    .menu-item:hover, .active {
        border-left: 3px solid #4CBB17;
        font-weight: bold;
    }
</style>

{#if $menuOptionsRO.length > 0}
    <div id="sidebar" class="flex sticky top-28 content-start">
        <div class="bg-gray-100 p-4 w-64 rounded shadow-md hidden md:block" >
            <nav>
                {#each $menuOptionsRO as section(section.id)}
                    <a href={"#" + section.id}
                       class="menu-item {section.id === activeSection ? 'active' : ''}">
                        {section.label}
                    </a>
                {/each}
            </nav>
        </div>
    </div>
{/if}
