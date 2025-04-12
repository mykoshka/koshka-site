<script lang="ts">
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { shoppingCart } from '$lib/stores/shopingCart'

    // In a real app, this would be a store that persists cart data
    let cartItems = [
        {
            id: 1,
            name: 'Organic Grain-Free Kibble',
            price: 34.99,
            quantity: 1,
            image: '/images/products/organic-kibble.jpg'
        },
        {
            id: 3,
            name: 'Modern Cat Tree',
            price: 119.99,
            quantity: 1,
            image: '/images/products/cat-tree.jpg'
        }
    ];

    // Computed values
    $: subtotal = $shoppingCart.reduce((sum, item) => sum + (100 * item.quantity), 0);
    $: shipping = subtotal > 100 ? 0 : 5.99;
    $: tax = subtotal * 0.15; // 15% tax rate
    $: total = subtotal + shipping + tax;

    // Update quantity
    function updateQuantity(index:number, newQuantity:number) {
        if (newQuantity < 1) {
            if (confirm('Remove this item from your cart?')) {
                $shoppingCart = $shoppingCart.filter((_, i) => i !== index);
            } else {
                $shoppingCart[index].quantity = 1;
                $shoppingCart = [...$shoppingCart]; // Trigger reactivity
            }
        } else {
            $shoppingCart[index].quantity = newQuantity;
            $shoppingCart = [...$shoppingCart]; // Trigger reactivity
        }
    }

    // Remove item
    function removeItem(index) {
        if (confirm('Are you sure you want to remove this item?')) {
            $shoppingCart = $shoppingCart.filter((_, i) => i !== index);
        }
    }

    // Proceed to checkout
    function proceedToCheckout() {
        // In a real app, this would save the cart state
        goto('/checkout');
    }

    // Continue shopping
    function continueShopping() {
        goto('/shop');
    }

    // For SEO
    onMount(() => {
        document.title = 'Your Cart | Koahka';
    });
</script>

<div class="bg-gray-50 min-h-screen py-8">
    <div class="container mx-auto px-4">
        <h1 class="text-3xl font-bold mb-8">Your Cart</h1>

        {#if $shoppingCart.length > 0}
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <!-- Cart Items -->
                <div class="lg:col-span-2">
                    <div class="bg-white rounded-lg shadow-sm overflow-hidden">
                        <!-- Header -->
                        <div class="grid grid-cols-12 gap-4 p-4 border-b bg-gray-50 hidden md:grid">
                            <div class="col-span-6">
                                <h2 class="font-medium text-gray-700">Product</h2>
                            </div>
                            <div class="col-span-2 text-center">
                                <h2 class="font-medium text-gray-700">Price</h2>
                            </div>
                            <div class="col-span-2 text-center">
                                <h2 class="font-medium text-gray-700">Quantity</h2>
                            </div>
                            <div class="col-span-2 text-right">
                                <h2 class="font-medium text-gray-700">Total</h2>
                            </div>
                        </div>

                        <!-- Items -->
                        {#each $shoppingCart as item, index}
                            <div class="grid grid-cols-1 md:grid-cols-12 gap-4 p-4 border-b items-center">
                                <!-- Product -->
                                <div class="md:col-span-6 flex items-center gap-4">
                                    <img
                                            src='/images/products/cat-tree.jpg'
                                            alt='Modern Cat Tree'
                                            class="w-16 h-16 object-cover rounded"
                                    />
                                    <div>
                                        <h3 class="font-medium text-gray-800">'Modern Cat Tree'</h3>
                                        <button
                                                on:click={() => removeItem(index)}
                                                class="text-sm text-red-500 hover:text-red-700 mt-1"
                                        >
                                            Remove
                                        </button>
                                    </div>
                                </div>

                                <!-- Price -->
                                <div class="md:col-span-2 text-center">
                                    <div class="md:hidden inline-block font-medium text-gray-700 mr-2">Price:</div>
                                    <span class="text-gray-800">${100.00}</span>
                                </div>

                                <!-- Quantity -->
                                <div class="md:col-span-2 flex justify-center">
                                    <div class="flex border rounded-md">
                                        <button
                                                on:click={() => updateQuantity(index, item.quantity - 1)}
                                                class="px-3 py-1 bg-gray-100 hover:bg-gray-200 transition-colors"
                                        >-</button>
                                        <input
                                                type="text"
                                                value={item.quantity}
                                                on:change={(e) => updateQuantity(index, parseInt(e.target.value))}
                                                min="1"
                                                class="w-12 text-center border-x py-1"
                                        />
                                        <button
                                                on:click={() => updateQuantity(index, item.quantity + 1)}
                                                class="px-3 py-1 bg-gray-100 hover:bg-gray-200 transition-colors"
                                        >+</button>
                                    </div>
                                </div>

                                <!-- Total -->
                                <div class="md:col-span-2 text-right">
                                    <div class="md:hidden inline-block font-medium text-gray-700 mr-2">Total:</div>
                                    <span class="font-medium text-gray-800">${(100 * item.quantity).toFixed(2)}</span>
                                </div>
                            </div>
                        {/each}
                    </div>

                    <!-- Continue Shopping -->
                    <div class="mt-6">
                        <button
                                on:click={continueShopping}
                                class="flex items-center text-emerald-600 hover:text-emerald-700"
                        >
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                            </svg>
                            Continue Shopping
                        </button>
                    </div>
                </div>

                <!-- Order Summary -->
                <div class="lg:col-span-1">
                    <div class="bg-white rounded-lg shadow-sm p-6">
                        <h2 class="text-xl font-bold mb-4">Order Summary</h2>

                        <div class="space-y-3 mb-6">
                            <div class="flex justify-between">
                                <span class="text-gray-600">Subtotal</span>
                                <span class="font-medium">${subtotal.toFixed(2)}</span>
                            </div>
                            <div class="flex justify-between">
                                <span class="text-gray-600">Shipping</span>
                                <span class="font-medium">
                                    {shipping === 0 ? 'Free' : `$${shipping.toFixed(2)}`}
                                </span>
                            </div>
                            <div class="flex justify-between">
                                <span class="text-gray-600">Tax</span>
                                <span class="font-medium">${tax.toFixed(2)}</span>
                            </div>
                            <div class="border-t pt-3 mt-3">
                                <div class="flex justify-between font-bold">
                                    <span>Total</span>
                                    <span>${total.toFixed(2)}</span>
                                </div>
                            </div>
                        </div>

                        <!-- Checkout Button -->
                        <button
                                on:click={proceedToCheckout}
                                class="w-full bg-emerald-500 hover:bg-emerald-600 text-white py-3 rounded-md transition-colors"
                        >
                            Proceed to Checkout
                        </button>

                        <!-- Shipping Note -->
                        <p class="text-sm text-gray-500 mt-4 text-center">
                            {subtotal >= 100
                                ? 'Your order qualifies for free shipping!'
                                : `Add $${(100 - subtotal).toFixed(2)} more to qualify for free shipping.`}
                        </p>
                    </div>
                </div>
            </div>
        {:else}
            <!-- Empty Cart -->
            <div class="bg-white rounded-lg p-8 text-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
                <h2 class="text-2xl font-bold text-gray-700 mb-2">Your cart is empty</h2>
                <p class="text-gray-500 mb-6">Looks like you haven't added any products to your cart yet.</p>
                <button
                        on:click={continueShopping}
                        class="bg-emerald-500 hover:bg-emerald-600 text-white px-6 py-3 rounded-md transition-colors"
                >
                    Start Shopping
                </button>
            </div>
        {/if}
    </div>
</div>
