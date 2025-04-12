<script>
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { addCartItem, shoppingCart } from '$lib/stores/shopingCart.js'

    // In a real app, this would be fetched from an API
    const products = [
        {
            id: 1,
            name: 'Organic Grain-Free Kibble',
            description: 'For sensitive digestive systems',
            price: 34.99,
            category: 'food',
            image: '/images/products/organic-kibble.jpg',
            longDescription: 'Our premium organic grain-free kibble is specially formulated for cats with sensitive digestive systems. Made with high-quality protein sources and free from common allergens, this food helps promote healthy digestion and overall wellbeing. Each batch is carefully prepared with natural ingredients and essential nutrients to keep your feline friend happy and healthy.',
            features: [
                'Grain-free formula for sensitive stomachs',
                'High-quality protein from organic sources',
                'No artificial preservatives or flavors',
                'Added vitamins and minerals for complete nutrition',
                'Supports digestive health and immune system'
            ],
            specifications: {
                weight: '4 lbs (1.8 kg)',
                dimensions: '10" x 5" x 15"',
                ingredients: 'Organic chicken, sweet potatoes, peas, natural flavors, fish oil, vitamins and minerals',
                feeding: 'Feed 1/2 to 3/4 cup daily for average adult cats'
            },
            reviews: [
                { author: 'Cat Parent', rating: 5, comment: 'My cat loves this food and her digestive issues have improved!' },
                { author: 'Fluffy\'s Owner', rating: 4, comment: 'Good quality food, my picky eater actually eats it.' }
            ]
        },
        {
            id: 2,
            name: 'Interactive Laser Toy',
            description: 'Motion activated for solo play',
            price: 24.99,
            category: 'toys',
            image: '/images/products/laser-toy.jpg',
            longDescription: 'Keep your cat entertained for hours with this interactive laser toy. The motion-activated design creates unpredictable patterns that trigger your cat\'s hunting instincts. Perfect for keeping indoor cats active and engaged, even when you\'re not home. The automatic timer ensures your cat gets play sessions throughout the day without constant supervision.',
            features: [
                'Motion-activated laser patterns',
                'Adjustable speed settings',
                'Automatic timer with 15, 30, and 60-minute sessions',
                'Battery-powered (3 AAA batteries included)',
                'Safe, low-power laser that won\'t harm eyes'
            ],
            specifications: {
                weight: '0.5 lbs (0.23 kg)',
                dimensions: '4" x 4" x 6"',
                batteryLife: 'Up to 8 hours of continuous play',
                materials: 'BPA-free plastic, metal components'
            },
            reviews: [
                { author: 'PlayfulCat', rating: 5, comment: 'My cat is obsessed with this toy! Great for keeping her active.' },
                { author: 'KittyMom', rating: 3, comment: 'Works well but battery life could be better.' }
            ]
        },
        {
            id: 3,
            name: 'Modern Cat Tree',
            description: 'Sleek design with sisal scratching posts',
            price: 119.99,
            category: 'beds',
            image: '/images/products/cat-tree.jpg',
            longDescription: 'This modern cat tree combines functionality with contemporary design to complement your home décor. Multiple platforms at different heights allow your cat to climb and perch, while the integrated sisal scratching posts help maintain healthy claws and satisfy natural scratching instincts. The plush-lined hideaway provides a cozy retreat for naps and relaxation.',
            features: [
                'Multiple platforms at varying heights',
                'Sisal-wrapped scratching posts',
                'Plush-lined hideaway cave',
                'Sturdy base for stability',
                'Easy to assemble with included tools'
            ],
            specifications: {
                weight: '25 lbs (11.3 kg)',
                dimensions: '28" x 24" x 48"',
                materials: 'Engineered wood, sisal rope, plush fabric',
                maxWeight: 'Supports cats up to 20 lbs'
            },
            reviews: [
                { author: 'ModernCatLover', rating: 5, comment: 'Beautiful design that actually looks good in my living room!' },
                { author: 'MultiCatHousehold', rating: 4, comment: 'Our three cats all love this tree. Sturdy and well-made.' }
            ]
        },
        {
            id: 4,
            name: 'Plush Donut Bed',
            description: 'Extra soft with calming texture',
            price: 49.99,
            category: 'beds',
            image: '/images/products/donut-bed.jpg',
            longDescription: 'Give your cat the ultimate comfort with our plush donut bed. The raised rim creates a sense of security and provides head and neck support, while the super-soft filling offers joint and muscle pain relief. The self-warming fabric retains your cat\'s body heat for extra coziness during naps. The non-slip bottom keeps the bed in place even during active moments.',
            features: [
                'Raised rim for security and support',
                'Ultra-plush filling for joint support',
                'Self-warming fabric technology',
                'Non-slip bottom',
                'Machine washable for easy cleaning'
            ],
            specifications: {
                weight: '3 lbs (1.4 kg)',
                dimensions: '20" diameter x 6" height',
                materials: 'Faux fur, recycled plush filling',
                care: 'Machine washable, tumble dry low'
            },
            reviews: [
                { author: 'SleepyCat', rating: 5, comment: 'My cat hasn\'t left this bed since it arrived! So soft and cozy.' },
                { author: 'SeniorCatOwner', rating: 5, comment: 'Perfect for my older cat with arthritis. She seems much more comfortable.' }
            ]
        },
        {
            id: 5,
            name: 'Dental Health Treats',
            description: 'Promotes healthy teeth and gums',
            price: 12.99,
            category: 'health',
            image: '/images/products/organic-kibble.jpg',
            longDescription: 'Our dental health treats are specially formulated to promote oral hygiene in cats. The crunchy texture helps reduce plaque and tartar buildup, while the natural ingredients freshen breath and support overall dental health. These treats are low in calories, making them perfect for daily use without contributing to weight gain.',
            features: [
                'Crunchy texture for plaque reduction',
                'Natural ingredients for fresh breath',
                'Low calorie formula',
                'No artificial colors or preservatives',
                'Supports overall dental health'
            ],
            specifications: {
                weight: '3 oz (85g)',
                dimensions: '5" x 3" x 8" package',
                ingredients: 'Chicken meal, brown rice, natural flavor, parsley (for breath freshening)',
                feeding: '5-10 treats daily for adult cats'
            },
            reviews: [
                { author: 'VetTech', rating: 4, comment: 'Good supplement to regular dental care. My cats enjoy them.' },
                { author: 'PickyCatMom', rating: 5, comment: 'Finally found dental treats my cat will actually eat!' }
            ]
        },
        {
            id: 6,
            name: 'Catnip Mouse Toy',
            description: 'Filled with premium catnip',
            price: 8.99,
            category: 'toys',
            image: '/images/products/laser-toy.jpg',
            longDescription: 'This adorable mouse toy is filled with premium, high-potency catnip to entice even the most discerning felines. The durable fabric exterior can withstand enthusiastic play, while the realistic size and texture appeal to your cat\'s natural hunting instincts. Each toy is handcrafted with reinforced stitching for longevity.',
            features: [
                'Filled with premium North American catnip',
                'Durable fabric exterior',
                'Realistic size and texture',
                'Reinforced stitching for longevity',
                'No small parts that could be swallowed'
            ],
            specifications: {
                weight: '0.1 lbs (45g)',
                dimensions: '3" x 1.5" x 1"',
                materials: 'Cotton fabric, premium catnip',
                care: 'Spot clean only'
            },
            reviews: [
                { author: 'CatnipLover', rating: 5, comment: 'My cat goes crazy for this toy! The catnip seems very potent.' },
                { author: 'ToyCollector', rating: 4, comment: 'Good quality and lasts longer than other similar toys.' }
            ]
        },
        {
            id: 7,
            name: 'Adjustable Collar with Bell',
            description: 'Soft fabric with safety release',
            price: 14.99,
            category: 'accessories',
            image: '/images/products/cat-tree.jpg',
            longDescription: 'Our adjustable collar features soft, comfortable fabric that\'s gentle on your cat\'s fur and skin. The breakaway safety buckle releases under pressure to prevent accidents if the collar gets caught. The included bell helps you keep track of your cat\'s whereabouts, and the collar is fully adjustable to ensure a perfect fit for cats of all sizes.',
            features: [
                'Soft, comfortable fabric material',
                'Breakaway safety buckle',
                'Removable bell',
                'Fully adjustable for perfect fit',
                'Reflective stitching for nighttime visibility'
            ],
            specifications: {
                weight: '0.1 lbs (45g)',
                dimensions: 'Adjusts from 8" to 12" length, 0.5" width',
                materials: 'Nylon fabric, plastic safety buckle',
                care: 'Hand wash with mild soap, air dry'
            },
            reviews: [
                { author: 'SafetyConcious', rating: 5, comment: 'Love the breakaway feature. Gives me peace of mind.' },
                { author: 'OutdoorCatOwner', rating: 4, comment: 'Good quality and the reflective stitching is helpful at night.' }
            ]
        },
        {
            id: 8,
            name: 'Window Perch',
            description: 'Strong suction cups for window viewing',
            price: 29.99,
            category: 'beds',
            image: '/images/products/donut-bed.jpg',
            longDescription: 'Give your cat the perfect vantage point with our window perch. The strong suction cups attach securely to any clean window, providing a comfortable spot for your cat to watch the outside world. The removable cover is machine washable for easy cleaning, and the sturdy platform supports cats up to 40 pounds. Installation is quick and tool-free.',
            features: [
                'Industrial-strength suction cups',
                'Comfortable padded surface',
                'Removable, washable cover',
                'Tool-free installation',
                'Supports cats up to 40 pounds'
            ],
            specifications: {
                weight: '2 lbs (0.9 kg)',
                dimensions: '12" x 22" x 2"',
                materials: 'Steel frame, oxford cloth cover, foam padding',
                care: 'Cover is machine washable'
            },
            reviews: [
                { author: 'WindowWatcher', rating: 5, comment: 'My cat spends hours on this perch! Suction cups are very strong.' },
                { author: 'ApartmentDweller', rating: 4, comment: 'Great for cats in small spaces. Easy to install and remove.' }
            ]
        }
    ];

    // Get product ID from URL parameter
    let productId = parseInt($page.params.id);

    // Find the product
    let product = products.find(p => p.id === productId);

    // Related products (in a real app, this would be more sophisticated)
    $: relatedProducts = products
        .filter(p => p.category === product?.category && p.id !== product?.id)
        .slice(0, 4);

    // Quantity for add to cart
    let quantity = 1;

    // For SEO
    onMount(() => {
        if (product) {
            document.title = `${product.name} | Koahka`;
        }
    });

    // Add to cart function
    function addToCart() {
        // In a real app, this would add to a cart store or send to an API
        addCartItem({prodID: productId.toString(), quantity: quantity})
        alert(`Added ${quantity} ${product.name} to cart!`);
        console.log($shoppingCart)
    }
</script>

{#if product}
    <div class="bg-white min-h-screen">
        <div class="container mx-auto px-4 py-8">
            <!-- Breadcrumbs -->
            <div class="text-sm text-gray-500 mb-6">
                <a href="/" class="hover:text-emerald-500">Home</a> &gt;
                <a href="/shop" class="hover:text-emerald-500">Shop</a> &gt;
                <span class="text-gray-700">{product.name}</span>
            </div>

            <!-- Product Details -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-8 mb-12">
                <!-- Product Image -->
                <div class="bg-gray-50 rounded-lg overflow-hidden">
                    <img
                            src={product.image}
                            alt={product.name}
                            class="w-full h-auto object-cover"
                    />
                </div>

                <!-- Product Info -->
                <div>
                    <h1 class="text-3xl font-bold text-gray-800 mb-2">{product.name}</h1>
                    <p class="text-gray-600 mb-4">{product.description}</p>

                    <div class="text-2xl font-bold text-emerald-600 mb-6">${product.price}</div>

                    <!-- Add to Cart -->
                    <div class="flex items-center gap-4 mb-8">
                        <div class="flex border rounded-md">
                            <button
                                    on:click={() => quantity = Math.max(1, quantity - 1)}
                                    class="px-3 py-2 bg-gray-100 hover:bg-gray-200 transition-colors"
                            >-</button>
                            <input
                                    type="text"
                                    bind:value={quantity}
                                    min="1"
                                    class="w-16 text-center border-x py-2"
                            />
                            <button
                                    on:click={() => quantity += 1}
                                    class="px-3 py-2 bg-gray-100 hover:bg-gray-200 transition-colors"
                            >+</button>
                        </div>

                        <button
                                on:click={addToCart}
                                class="bg-emerald-500 hover:bg-emerald-600 text-white px-6 py-2 rounded-md transition-colors flex-grow md:flex-grow-0"
                        >
                            Add to Cart
                        </button>
                    </div>

                    <!-- Product Description -->
                    <div class="mb-6">
                        <h2 class="text-xl font-semibold mb-2">Description</h2>
                        <p class="text-gray-700">{product.longDescription}</p>
                    </div>

                    <!-- Product Features -->
                    <div class="mb-6">
                        <h2 class="text-xl font-semibold mb-2">Features</h2>
                        <ul class="list-disc pl-5 text-gray-700">
                            {#each product.features as feature}
                                <li class="mb-1">{feature}</li>
                            {/each}
                        </ul>
                    </div>

                    <!-- Product Specifications -->
                    <div>
                        <h2 class="text-xl font-semibold mb-2">Specifications</h2>
                        <div class="bg-gray-50 p-4 rounded-lg">
                            {#each Object.entries(product.specifications) as [key, value]}
                                <div class="grid grid-cols-3 py-2 border-b border-gray-200 last:border-0">
                                    <div class="font-medium text-gray-700 capitalize">{key}</div>
                                    <div class="col-span-2 text-gray-600">{value}</div>
                                </div>
                            {/each}
                        </div>
                    </div>
                </div>
            </div>

            <!-- Reviews Section -->
            <div class="mb-12">
                <h2 class="text-2xl font-bold mb-4">Customer Reviews</h2>

                {#if product.reviews && product.reviews.length > 0}
                    <div class="space-y-4">
                        {#each product.reviews as review}
                            <div class="bg-gray-50 p-4 rounded-lg">
                                <div class="flex items-center mb-2">
                                    <div class="flex text-yellow-400 mr-2">
                                        {#each Array(5) as _, i}
                                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 {i < review.rating ? 'text-yellow-400' : 'text-gray-300'}" viewBox="0 0 20 20" fill="currentColor">
                                                <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                                            </svg>
                                        {/each}
                                    </div>
                                    <span class="font-medium text-gray-700">{review.author}</span>
                                </div>
                                <p class="text-gray-600">{review.comment}</p>
                            </div>
                        {/each}
                    </div>
                {:else}
                    <p class="text-gray-600">No reviews yet. Be the first to review this product!</p>
                {/if}
            </div>

            <!-- Related Products -->
            {#if relatedProducts.length > 0}
                <div>
                    <h2 class="text-2xl font-bold mb-4">You May Also Like</h2>

                    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
                        {#each relatedProducts as relatedProduct}
                            <a href="/shop/{relatedProduct.id}" class="bg-white rounded-lg overflow-hidden shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
                                <div class="aspect-square bg-gray-50 relative overflow-hidden">
                                    <img
                                            src={relatedProduct.image}
                                            alt={relatedProduct.name}
                                            class="w-full h-full object-cover"
                                            loading="lazy"
                                    />
                                </div>
                                <div class="p-4">
                                    <h3 class="font-medium text-gray-800 mb-1">{relatedProduct.name}</h3>
                                    <p class="text-gray-600 text-sm mb-3">{relatedProduct.description}</p>
                                    <div class="flex justify-between items-center">
                                        <span class="font-medium text-emerald-600">${relatedProduct.price}</span>
                                        <button class="bg-emerald-500 text-white text-sm py-1.5 px-3 rounded hover:bg-emerald-600 transition-colors">
                                            View
                                        </button>
                                    </div>
                                </div>
                            </a>
                        {/each}
                    </div>
                </div>
            {/if}
        </div>
    </div>
{:else}
    <div class="bg-white min-h-screen flex items-center justify-center">
        <div class="text-center p-8">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <h1 class="text-2xl font-bold text-gray-700 mb-2">Product Not Found</h1>
            <p class="text-gray-500 mb-4">The product you're looking for doesn't exist or has been removed.</p>
            <a href="/shop" class="inline-block bg-emerald-500 text-white px-6 py-2 rounded-md hover:bg-emerald-600 transition-colors">
                Return to Shop
            </a>
        </div>
    </div>
{/if}
