<script lang="ts">
    import type {IuserProfile} from '$lib/stores/authStore'
    import {userProfile} from '$lib/stores/authStore'
    let userProfileValue: IuserProfile;

    userProfile.subscribe(value => {
        userProfileValue = value;
    });

</script>

<!-- Header -->
<div class="sticky top-0 z-50">
<header class="relative bg-white shadow-sm">
    <div class="container mx-auto px-4 py-4 flex justify-between items-center">
        <!-- Logo and Brand -->
        <a href="/" class="flex items-center gap-2" aria-label="Koahka Home">
            <enhanced:img src="/static/images/company-logo.png" alt="An alt text" />
        </a>

        <!-- Navigation -->
        <nav class="hidden md:flex items-center space-x-6">
            <a href="/products/collars" class="text-gray-600 hover:text-emerald-500 transition-colors">Collars</a>
            <a href="/products/scratchposts" class="text-gray-600 hover:text-emerald-500 transition-colors">Scratch Posts</a>
            <a href="/blog" class="text-gray-600 hover:text-emerald-500 transition-colors">Blog</a>
            <a href="/products/instructions" class="text-gray-600 hover:text-emerald-500 transition-colors">Instructions</a>
        </nav>

        <!-- Search -->
        <div class="flex items-center gap-4">
            <div class="relative hidden md:block">
                <input
                        type="text"
                        placeholder="Search..."
                        class="bg-gray-100 rounded-full py-2 px-4 pr-10 w-64 focus:outline-none focus:ring-2 focus:ring-emerald-300"
                />
                <button class="absolute right-3 top-2.5" aria-label="Search">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                </button>
            </div>

            <!-- User Icons -->
            <button class="p-2" aria-label="Notifications">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                </svg>
            </button>

            <a href={ userProfileValue.isLoggedIn ? "/profile" : "/login"}>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 {userProfileValue.isLoggedIn ? 'text-green-600' : 'text-gray-700'}" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
            </a>

            <!-- Mobile Menu Button -->
            <button class="p-2 md:hidden" aria-label="Menu">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                </svg>
            </button>
        </div>
    </div>
</header>
</div>