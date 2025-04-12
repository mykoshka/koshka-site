<script>
    import { onMount } from 'svelte';

    // Product categories
    const categories = [
        { id: 'food', name: 'Food', icon: 'bowl' },
        { id: 'toys', name: 'Toys', icon: 'toy' },
        { id: 'beds', name: 'Beds', icon: 'bed' },
        { id: 'health', name: 'Health', icon: 'health' },
        { id: 'accessories', name: 'Accessories', icon: 'collar' }
    ];

    // Products data - in a real app, this would come from an API
    const products = [
        {
            id: 1,
            name: 'Organic Grain-Free Kibble',
            description: 'For sensitive digestive systems',
            price: 34.99,
            category: 'food',
            image: '/images/products/organic-kibble.jpg'
        },
        {
            id: 2,
            name: 'Interactive Laser Toy',
            description: 'Motion activated for solo play',
            price: 24.99,
            category: 'toys',
            image: '/images/products/laser-toy.jpg'
        },
        {
            id: 3,
            name: 'Modern Cat Tree',
            description: 'Sleek design with sisal scratching posts',
            price: 119.99,
            category: 'beds',
            image: '/images/products/cat-tree.jpg'
        },
        {
            id: 4,
            name: 'Plush Donut Bed',
            description: 'Extra soft with calming texture',
            price: 49.99,
            category: 'beds',
            image: '/images/products/donut-bed.jpg'
        },
        {
            id: 5,
            name: 'Dental Health Treats',
            description: 'Promotes healthy teeth and gums',
            price: 12.99,
            category: 'health',
            image: '/images/products/organic-kibble.jpg'
        },
        {
            id: 6,
            name: 'Catnip Mouse Toy',
            description: 'Filled with premium catnip',
            price: 8.99,
            category: 'toys',
            image: '/images/products/laser-toy.jpg'
        },
        {
            id: 7,
            name: 'Adjustable Collar with Bell',
            description: 'Soft fabric with safety release',
            price: 14.99,
            category: 'accessories',
            image: '/images/products/cat-tree.jpg'
        },
        {
            id: 8,
            name: 'Window Perch',
            description: 'Strong suction cups for window viewing',
            price: 29.99,
            category: 'beds',
            image: '/images/products/donut-bed.jpg'
        }
    ];

    // Filter state
    let selectedCategory = '';
    let searchQuery = '';
    let sortOption = 'featured';

    // Computed filtered products
    $: filteredProducts = products.filter(product => {
        // Filter by category if selected
        if (selectedCategory && product.category !== selectedCategory) {
            return false;
        }

        // Filter by search query
        if (searchQuery && !product.name.toLowerCase().includes(searchQuery.toLowerCase())) {
            return false;
        }

        return true;
    }).sort((a, b) => {
        // Sort products based on selected option
        if (sortOption === 'price-low') {
            return a.price - b.price;
        } else if (sortOption === 'price-high') {
            return b.price - a.price;
        } else if (sortOption === 'name') {
            return a.name.localeCompare(b.name);
        }
        // Default: featured (no specific sort)
        return 0;
    });

    // For SEO
    onMount(() => {
        document.title = 'Shop Cat Products | Koahka';
    });

    // Handle category selection
    function selectCategory(categoryId) {
        selectedCategory = categoryId === selectedCategory ? '' : categoryId;
    }
</script>

<div class="bg-gray-50 min-h-screen py-8">
    <div class="container mx-auto px-4">
        <h1 class="text-3xl font-bold mb-8">Shop Cat Products</h1>

        <!-- Filters and Search -->
        <div class="bg-white p-4 rounded-lg shadow-sm mb-8">
            <div class="flex flex-col md:flex-row md:justify-between md:items-center gap-4">
                <!-- Search -->
                <div class="relative w-full md:w-64">
                    <input
                            type="text"
                            bind:value={searchQuery}
                            placeholder="Search products..."
                            class="w-full py-2 px-4 pr-10 border rounded-full focus:outline-none focus:ring-2 focus:ring-emerald-300"
                    />
                    <span class="absolute right-3 top-2.5">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                        </svg>
                    </span>
                </div>

                <!-- Sort Options -->
                <div class="flex items-center">
                    <label for="sort" class="text-sm text-gray-600 mr-2">Sort by:</label>
                    <select
                            id="sort"
                            bind:value={sortOption}
                            class="border rounded-lg py-2 px-3 focus:outline-none focus:ring-2 focus:ring-emerald-300"
                    >
                        <option value="featured">Featured</option>
                        <option value="price-low">Price: Low to High</option>
                        <option value="price-high">Price: High to Low</option>
                        <option value="name">Name</option>
                    </select>
                </div>
            </div>

            <!-- Category Filters -->
            <div class="mt-4 flex flex-wrap gap-2">
                {#each categories as category}
                    <button
                            on:click={() => selectCategory(category.id)}
                            class="px-4 py-2 rounded-full text-sm {selectedCategory === category.id ? 'bg-emerald-500 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'} transition-colors"
                    >
                        {category.name}
                    </button>
                {/each}
                {#if selectedCategory}
                    <button
                            on:click={() => selectedCategory = ''}
                            class="px-4 py-2 rounded-full text-sm bg-gray-200 text-gray-700 hover:bg-gray-300 transition-colors"
                    >
                        Clear Filter
                    </button>
                {/if}
            </div>
        </div>

        <!-- Product Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {#each filteredProducts as product}
                <a href="/shop/{product.id}" class="bg-white rounded-lg overflow-hidden shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
                    <div class="aspect-square bg-gray-50 relative overflow-hidden">
                        <img
                                src={product.image}
                                alt={product.name}
                                class="w-full h-full object-cover"
                                loading="lazy"
                        />
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
                </a>
            {/each}
        </div>

        <!-- Empty State -->
        {#if filteredProducts.length === 0}
            <div class="bg-white rounded-lg p-8 text-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <h3 class="text-lg font-medium text-gray-700 mb-2">No products found</h3>
                <p class="text-gray-500">Try adjusting your search or filter to find what you're looking for.</p>
                <button
                        on:click={() => { searchQuery = ''; selectedCategory = ''; }}
                        class="mt-4 px-4 py-2 bg-emerald-500 text-white rounded-lg hover:bg-emerald-600 transition-colors"
                >
                    Clear All Filters
                </button>
            </div>
        {/if}
    </div>
</div>
