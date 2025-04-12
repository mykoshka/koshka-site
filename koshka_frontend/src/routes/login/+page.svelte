<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import gsap from "gsap";
    import axios from "axios";
    import {storeTokens} from "$lib/auth";

    let email = "";
    let password = "";
    let errorMessage = "";

    async function login() {
        try {
            const response = await axios.post("http://localhost:3000/login", { email, password });
            storeTokens(response.data.tokens.auth_token, response.data.tokens.refresh_token);
            goto("/profile"); // Redirect to profile after login
        } catch {
            errorMessage = "Invalid email or password";
        }
    }

    onMount(() => {
        gsap.from(".login-container", { opacity: 0, y: -50, duration: 1 });
    });
</script>

<div class="flex items-center justify-center min-h-screen bg-gray-100">
    <div class="login-container bg-white p-8 rounded-lg shadow-lg w-96">
        <h2 class="text-2xl font-semibold text-center">Login</h2>

        {#if errorMessage}
            <p class="text-red-500 text-sm mt-2">{errorMessage}</p>
        {/if}

        <input
                type="email"
                bind:value={email}
                placeholder="Email"
                class="w-full p-2 mt-4 border rounded-lg"
        />

        <input
                type="password"
                bind:value={password}
                placeholder="Password"
                class="w-full p-2 mt-2 border rounded-lg"
        />

        <button
                on:click={login}
                class="w-full mt-4 bg-green-500 text-white py-2 rounded-lg hover:bg-green-600"
        >
            Login
        </button>
    </div>
</div>