<script lang="ts">
	import '../app.css';
    import { userProfile } from '$lib/stores/authStore'
    import { onMount } from 'svelte';
    import { isAuthenticated }  from '$lib/auth'
    import Header from  '$lib/components/Header.svelte'
    import Footer from  '$lib/components/Footer.svelte'

    onMount(async () => {
        const isauth = isAuthenticated()
        $userProfile = { isLoggedIn: isauth, role: null }
    })

</script>

<svelte:head>
    <!-- Additional SEO meta tags -->
    <meta name="keywords" content="cat products, premium cat food, cat toys, cat beds, feline health, cat accessories" />
    <meta property="og:title" content="Koahka - Premium Products for Cat Companions" />
    <meta property="og:description" content="Thoughtfully crafted premium products for the discerning feline." />
    <meta property="og:image" content="/images/koahka-banner.jpg" />
    <meta property="og:url" content="https://koahka.com" />
    <meta property="og:type" content="website" />
    <meta name="twitter:card" content="summary_large_image" />
    <!-- Canonical URL for SEO -->
    <link rel="canonical" href="https://koahka.com" />
</svelte:head>


<Header />
<div class="min-h-[75vh]">
    <slot />
</div>
<Footer />