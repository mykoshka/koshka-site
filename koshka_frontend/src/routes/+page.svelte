<!-- App.svelte -->
<script>
    import { onMount } from 'svelte';

    // Products data
    const bestSellers = [
        {
            id: 1,
            name: 'Organic Grain-Free Kibble',
            description: 'For sensitive digestive systems',
            price: 34.99,
            image: '/images/products/organic-kibble.jpg'
        },
        {
            id: 2,
            name: 'Interactive Laser Toy',
            description: 'Motion activated for solo play',
            price: 24.99,
            image: '/images/products/laser-toy.jpg'
        },
        {
            id: 3,
            name: 'Modern Cat Tree',
            description: 'Sleek design with sisal scratching posts',
            price: 119.99,
            image: '/images/products/cat-tree.jpg'
        },
        {
            id: 4,
            name: 'Plush Donut Bed',
            description: 'Extra soft with calming texture',
            price: 49.99,
            image: '/images/products/donut-bed.jpg'
        }
    ];

    // Categories
    const categories = [
        {
            name: 'Premium Food',
            icon: 'food-icon',
            link: '/food'
        },
        {
            name: 'Interactive Toys',
            icon: 'toys-icon',
            link: '/toys'
        },
        {
            name: 'Cozy Beds',
            icon: 'bed-icon',
            link: '/beds'
        },
        {
            name: 'Health Essentials',
            icon: 'health-icon',
            link: '/health'
        }
    ];

    // Email signup
    let email = '';
    const handleSubmit = () => {
        // Handle newsletter signup
        if (email) {
            alert('Thank you for joining our community!');
            email = '';
        }
    };

    // For SEO optimization
    onMount(() => {
        document.title = 'Koahka - Premium Products for Cat Companions';
        const metaDescription = document.createElement('meta');
        metaDescription.name = 'description';
        metaDescription.content = 'Thoughtfully crafted premium products for the discerning feline. Shop our selection of cat food, toys, beds, and health accessories.';
        document.head.appendChild(metaDescription);
    });
</script>

<main class="min-h-screen bg-white">

    <!-- Hero Section -->
    <section class="bg-emerald-50 py-16">
        <div class="container mx-auto px-4 md:flex md:items-center md:justify-between">
            <div class="md:w-1/2 mb-8 md:mb-0">
                <h1 class="text-4xl md:text-5xl font-bold text-gray-800 mb-2">
                    Thoughtfully Crafted for
                    <span class="text-emerald-500 block mt-2">Cat Companions</span>
                </h1>
                <p class="text-gray-600 mt-4 mb-8">Premium products for the discerning feline</p>
                <a
                        href="/shop"
                        class="bg-emerald-500 hover:bg-emerald-600 text-white font-medium py-3 px-6 rounded-full inline-block transition-colors uppercase text-sm"
                >
                    Shop Now
                </a>
            </div>
            <div class="md:w-1/2 flex justify-end">
                <!-- Cute cat face illustration -->
                <div class="w-64 h-64 md:w-80 md:h-80 relative">
                    <div class="absolute inset-0 bg-emerald-100 rounded-full opacity-50"></div>
                    <div class="absolute inset-8 bg-emerald-200 rounded-full flex items-center justify-center">
                        <!-- Cat face features -->
                        <div class="relative w-full h-full">
                            <div class="absolute w-2 h-2 bg-emerald-500 rounded-full left-1/3 top-1/3"></div>
                            <div class="absolute w-2 h-2 bg-emerald-500 rounded-full right-1/3 top-1/3"></div>
                            <div class="absolute w-10 h-4 bottom-1/3 left-1/2 transform -translate-x-1/2">
                                <svg viewBox="0 0 40 16" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M8 8C8 4 12 4 20 4C28 4 32 4 32 8" stroke="#4ade80" stroke-width="2" fill="none" stroke-linecap="round"/>
                                </svg>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>

    <!-- Categories Section -->
    <section class="py-12">
        <div class="container mx-auto px-4">
            <h2 class="text-2xl font-bold mb-8">Shop Categories</h2>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                {#each categories as category}
                    <a href={category.link} class="bg-white border border-gray-100 rounded-lg p-6 text-center hover:shadow-md transition-shadow">
                        <div class="w-24 h-24 bg-emerald-50 rounded-full mx-auto mb-4 flex items-center justify-center">
                            <!-- Category icon placeholders -->
                            {#if category.name === 'Premium Food'}
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-emerald-500" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                                    <circle cx="12" cy="12" r="10" stroke-width="1"/>
                                    <circle cx="9" cy="9" r="1.5" stroke-width="0.5"/>
                                    <circle cx="15" cy="9" r="1.5" stroke-width="0.5"/>
                                </svg>
                            {:else if category.name === 'Interactive Toys'}
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-emerald-500" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                                    <circle cx="12" cy="12" r="8" stroke-width="1"/>
                                    <line x1="4" y1="4" x2="20" y2="20" stroke-width="1"/>
                                    <line x1="4" y1="20" x2="20" y2="4" stroke-width="1"/>
                                </svg>
                            {:else if category.name === 'Cozy Beds'}
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-emerald-500" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                                    <path d="M3 17V7C3 5.89543 3.89543 5 5 5H19C20.1046 5 21 5.89543 21 7V17" stroke-width="1"/>
                                    <path d="M3 17H21" stroke-width="1"/>
                                </svg>
                            {:else if category.name === 'Health Essentials'}
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-emerald-500" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                                    <circle cx="12" cy="12" r="9" stroke-width="1"/>
                                    <line x1="12" y1="8" x2="12" y2="16" stroke-width="1"/>
                                    <line x1="8" y1="12" x2="16" y2="12" stroke-width="1"/>
                                </svg>
                            {/if}
                        </div>
                        <h3 class="font-medium text-gray-800">{category.name}</h3>
                    </a>
                {/each}
            </div>
        </div>
    </section>

    <!-- Best Sellers Section -->
    <section class="py-12 bg-gray-50">
        <div class="container mx-auto px-4">
            <h2 class="text-2xl font-bold mb-8">Best Sellers</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {#each bestSellers as product}
                    <div class="bg-white rounded-lg overflow-hidden shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
                        <div class="aspect-square bg-emerald-50 flex items-center justify-center">
                            <!-- Product image placeholders -->
                            <div class="w-32 h-32 bg-emerald-100 rounded-full"></div>
                        </div>
                        <div class="p-4">
                            <h3 class="font-medium text-gray-800 mb-1">{product.name}</h3>
                            <p class="text-gray-600 text-sm mb-3">{product.description}</p>
                            <div class="flex justify-between items-center">
                                <span class="font-medium text-emerald-600">${product.price}</span>
                                <button class="bg-emerald-500 text-white text-sm py-1.5 px-3 rounded hover:bg-emerald-600 transition-colors">
                                    Add to Cart
                                </button>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    </section>

    <!-- Newsletter Section -->
    <section class="py-12 bg-emerald-50">
        <div class="container mx-auto px-4">
            <div class="bg-white rounded-lg shadow-sm p-8 md:p-10">
                <div class="md:flex md:justify-between md:items-center">
                    <div class="mb-6 md:mb-0 md:w-1/2">
                        <h3 class="text-xl font-bold mb-2">Join our cat-loving community</h3>
                        <p class="text-gray-600">Get care tips and exclusive offers</p>
                    </div>
                    <div class="md:w-1/2">
                        <form on:submit|preventDefault={handleSubmit} class="flex">
                            <input
                                    type="email"
                                    bind:value={email}
                                    placeholder="Your email"
                                    required
                                    class="flex-grow rounded-l-full border-r-0 border border-gray-300 py-3 px-4 focus:outline-none focus:ring-2 focus:ring-emerald-300"
                            />
                            <button
                                    type="submit"
                                    class="bg-emerald-500 text-white px-6 py-3 rounded-r-full hover:bg-emerald-600 transition-colors uppercase font-medium"
                            >
                                Join
                            </button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </section>
</main>

<style>
    /* Any additional global styles could go here */
    :global(html) {
        scroll-behavior: smooth;
    }
</style>