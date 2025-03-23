<script>
    import { goto } from '$app/navigation';

    let email = '';
    let password = '';
    let name = '';
    let address = '';
    let mobile_number = '';
    let error = '';
    let success = '';

    async function handleRegister() {
        try {
            const response = await fetch('/api/auth/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password, name, address, mobile_number }),
            });
            const data = await response.json();
            if (response.ok) {
                success = data.message || 'Registration successful! Please log in.';
                setTimeout(() => goto('/login'), 2000);
            } else {
                error = data.message || 'Registration failed';
            }
        } catch (err) {
            error = 'An error occurred. Please try again.';
        }
    }
</script>

<div class="min-h-screen bg-koahka-bg flex items-center justify-center">
    <div class="bg-white p-8 rounded-lg shadow-lg w-full max-w-md">
        <h1 class="text-3xl font-bold text-center mb-6">Join Koahka</h1>
        {#if error}
            <p class="text-red-500 text-center mb-4">{error}</p>
        {/if}
        {#if success}
            <p class="text-koahka-green text-center mb-4">{success}</p>
        {/if}
        <form on:submit|preventDefault={handleRegister}>
            <div class="mb-4">
                <label for="name" class="block text-sm font-medium mb-2">Name</label>
                <input
                        type="text"
                        id="name"
                        bind:value={name}
                        class="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-koahka-green"
                        placeholder="Your name"
                        required
                />
            </div>
            <div class="mb-4">
                <label for="email" class="block text-sm font-medium mb-2">Email</label>
                <input
                        type="email"
                        id="email"
                        bind:value={email}
                        class="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-koahka-green"
                        placeholder="Your email"
                        required
                />
            </div>
            <div class="mb-4">
                <label for="password" class="block text-sm font-medium mb-2">Password</label>
                <input
                        type="password"
                        id="password"
                        bind:value={password}
                        class="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-koahka-green"
                        placeholder="Your password"
                        required
                />
            </div>
            <div class="mb-4">
                <label for="address" class="block text-sm font-medium mb-2">Address</label>
                <input
                        type="text"
                        id="address"
                        bind:value={address}
                        class="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-koahka-green"
                        placeholder="Your address"
                />
            </div>
            <div class="mb-6">
                <label for="mobile_number" class="block text-sm font-medium mb-2">Mobile Number</label>
                <input
                        type="text"
                        id="mobile_number"
                        bind:value={mobile_number}
                        class="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-koahka-green"
                        placeholder="Your mobile number"
                />
            </div>
            <button
                    type="submit"
                    class="w-full bg-koahka-green text-white py-3 rounded-full hover:bg-green-600 transition"
            >
                Register
            </button>
        </form>
        <p class="text-center mt-4 text-sm">
            Already have an account? <a href="/login" class="text-koahka-green hover:underline">Login</a>
        </p>
    </div>
</div>