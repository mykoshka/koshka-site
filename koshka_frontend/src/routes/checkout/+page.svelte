<script>
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';

    // Form data
    let formData = {
        firstName: '',
        lastName: '',
        email: '',
        address: '',
        city: '',
        state: '',
        zipCode: '',
        country: 'United States',
        phone: '',
        sameAsBilling: true,
        shippingAddress: '',
        shippingCity: '',
        shippingState: '',
        shippingZipCode: '',
        shippingCountry: 'United States',
        paymentMethod: 'credit',
        cardNumber: '',
        cardName: '',
        expiryDate: '',
        cvv: ''
    };

    // Validation
    let errors = {};
    let formSubmitted = false;

    // In a real app, this would be fetched from a cart store
    const cartItems = [
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
    $: subtotal = cartItems.reduce((sum, item) => sum + (item.price * item.quantity), 0);
    $: shipping = subtotal > 100 ? 0 : 5.99;
    $: tax = subtotal * 0.08; // 8% tax rate
    $: total = subtotal + shipping + tax;

    // Validate form
    function validateForm() {
        errors = {};

        // Required fields
        const requiredFields = ['firstName', 'lastName', 'email', 'address', 'city', 'state', 'zipCode', 'country', 'phone'];
        requiredFields.forEach(field => {
            if (!formData[field]) {
                errors[field] = 'This field is required';
            }
        });

        // Shipping address if not same as billing
        if (!formData.sameAsBilling) {
            const requiredShippingFields = ['shippingAddress', 'shippingCity', 'shippingState', 'shippingZipCode', 'shippingCountry'];
            requiredShippingFields.forEach(field => {
                if (!formData[field]) {
                    errors[field] = 'This field is required';
                }
            });
        }

        // Payment information
        if (formData.paymentMethod === 'credit') {
            if (!formData.cardNumber) errors.cardNumber = 'Card number is required';
            if (!formData.cardName) errors.cardName = 'Name on card is required';
            if (!formData.expiryDate) errors.expiryDate = 'Expiry date is required';
            if (!formData.cvv) errors.cvv = 'CVV is required';

            // Basic format validation
            if (formData.cardNumber && !/^\d{16}$/.test(formData.cardNumber.replace(/\s/g, ''))) {
                errors.cardNumber = 'Invalid card number format';
            }
            if (formData.expiryDate && !/^\d{2}\/\d{2}$/.test(formData.expiryDate)) {
                errors.expiryDate = 'Use MM/YY format';
            }
            if (formData.cvv && !/^\d{3,4}$/.test(formData.cvv)) {
                errors.cvv = 'Invalid CVV';
            }
        }

        // Email validation
        if (formData.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
            errors.email = 'Invalid email format';
        }

        return Object.keys(errors).length === 0;
    }

    // Handle form submission
    function handleSubmit() {
        formSubmitted = true;

        if (validateForm()) {
            // In a real app, this would submit the order to an API
            alert('Order placed successfully! Thank you for your purchase.');
            goto('/');
        } else {
            // Scroll to the first error
            const firstErrorField = document.querySelector('.error-field');
            if (firstErrorField) {
                firstErrorField.scrollIntoView({ behavior: 'smooth', block: 'center' });
            }
        }
    }

    // For SEO
    onMount(() => {
        document.title = 'Checkout | Koahka';
    });
</script>

<div class="bg-gray-50 min-h-screen py-8">
    <div class="container mx-auto px-4">
        <h1 class="text-3xl font-bold mb-8">Checkout</h1>

        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <!-- Checkout Form -->
            <div class="lg:col-span-2">
                <form on:submit|preventDefault={handleSubmit} class="bg-white rounded-lg shadow-sm p-6">
                    <!-- Contact Information -->
                    <div class="mb-8">
                        <h2 class="text-xl font-bold mb-4">Contact Information</h2>

                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div class="form-group">
                                <label for="firstName" class="block text-gray-700 mb-1">First Name *</label>
                                <input
                                        type="text"
                                        id="firstName"
                                        bind:value={formData.firstName}
                                        class="w-full p-2 border rounded-md {errors.firstName && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                        class:error-field={errors.firstName && formSubmitted}
                                />
                                {#if errors.firstName && formSubmitted}
                                    <p class="text-red-500 text-sm mt-1">{errors.firstName}</p>
                                {/if}
                            </div>

                            <div class="form-group">
                                <label for="lastName" class="block text-gray-700 mb-1">Last Name *</label>
                                <input
                                        type="text"
                                        id="lastName"
                                        bind:value={formData.lastName}
                                        class="w-full p-2 border rounded-md {errors.lastName && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                        class:error-field={errors.lastName && formSubmitted}
                                />
                                {#if errors.lastName && formSubmitted}
                                    <p class="text-red-500 text-sm mt-1">{errors.lastName}</p>
                                {/if}
                            </div>

                            <div class="form-group">
                                <label for="email" class="block text-gray-700 mb-1">Email *</label>
                                <input
                                        type="email"
                                        id="email"
                                        bind:value={formData.email}
                                        class="w-full p-2 border rounded-md {errors.email && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                        class:error-field={errors.email && formSubmitted}
                                />
                                {#if errors.email && formSubmitted}
                                    <p class="text-red-500 text-sm mt-1">{errors.email}</p>
                                {/if}
                            </div>

                            <div class="form-group">
                                <label for="phone" class="block text-gray-700 mb-1">Phone *</label>
                                <input
                                        type="tel"
                                        id="phone"
                                        bind:value={formData.phone}
                                        class="w-full p-2 border rounded-md {errors.phone && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                        class:error-field={errors.phone && formSubmitted}
                                />
                                {#if errors.phone && formSubmitted}
                                    <p class="text-red-500 text-sm mt-1">{errors.phone}</p>
                                {/if}
                            </div>
                        </div>
                    </div>

                    <!-- Billing Address -->
                    <div class="mb-8">
                        <h2 class="text-xl font-bold mb-4">Billing Address</h2>

                        <div class="grid grid-cols-1 gap-4">
                            <div class="form-group">
                                <label for="address" class="block text-gray-700 mb-1">Address *</label>
                                <input
                                        type="text"
                                        id="address"
                                        bind:value={formData.address}
                                        class="w-full p-2 border rounded-md {errors.address && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                        class:error-field={errors.address && formSubmitted}
                                />
                                {#if errors.address && formSubmitted}
                                    <p class="text-red-500 text-sm mt-1">{errors.address}</p>
                                {/if}
                            </div>

                            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                                <div class="form-group">
                                    <label for="city" class="block text-gray-700 mb-1">City *</label>
                                    <input
                                            type="text"
                                            id="city"
                                            bind:value={formData.city}
                                            class="w-full p-2 border rounded-md {errors.city && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                            class:error-field={errors.city && formSubmitted}
                                    />
                                    {#if errors.city && formSubmitted}
                                        <p class="text-red-500 text-sm mt-1">{errors.city}</p>
                                    {/if}
                                </div>

                                <div class="form-group">
                                    <label for="state" class="block text-gray-700 mb-1">State/Province *</label>
                                    <input
                                            type="text"
                                            id="state"
                                            bind:value={formData.state}
                                            class="w-full p-2 border rounded-md {errors.state && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                            class:error-field={errors.state && formSubmitted}
                                    />
                                    {#if errors.state && formSubmitted}
                                        <p class="text-red-500 text-sm mt-1">{errors.state}</p>
                                    {/if}
                                </div>

                                <div class="form-group">
                                    <label for="zipCode" class="block text-gray-700 mb-1">ZIP/Postal Code *</label>
                                    <input
                                            type="text"
                                            id="zipCode"
                                            bind:value={formData.zipCode}
                                            class="w-full p-2 border rounded-md {errors.zipCode && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                            class:error-field={errors.zipCode && formSubmitted}
                                    />
                                    {#if errors.zipCode && formSubmitted}
                                        <p class="text-red-500 text-sm mt-1">{errors.zipCode}</p>
                                    {/if}
                                </div>
                            </div>

                            <div class="form-group">
                                <label for="country" class="block text-gray-700 mb-1">Country *</label>
                                <select
                                        id="country"
                                        bind:value={formData.country}
                                        class="w-full p-2 border rounded-md {errors.country && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                        class:error-field={errors.country && formSubmitted}
                                >
                                    <option value="United States">United States</option>
                                    <option value="Canada">Canada</option>
                                    <option value="United Kingdom">United Kingdom</option>
                                    <option value="Australia">Australia</option>
                                    <option value="Germany">Germany</option>
                                    <option value="France">France</option>
                                </select>
                                {#if errors.country && formSubmitted}
                                    <p class="text-red-500 text-sm mt-1">{errors.country}</p>
                                {/if}
                            </div>
                        </div>
                    </div>

                    <!-- Shipping Address -->
                    <div class="mb-8">
                        <div class="flex items-center mb-4">
                            <input
                                    type="checkbox"
                                    id="sameAsBilling"
                                    bind:checked={formData.sameAsBilling}
                                    class="mr-2"
                            />
                            <label for="sameAsBilling" class="text-gray-700">Shipping address same as billing</label>
                        </div>

                        {#if !formData.sameAsBilling}
                            <h2 class="text-xl font-bold mb-4">Shipping Address</h2>

                            <div class="grid grid-cols-1 gap-4">
                                <div class="form-group">
                                    <label for="shippingAddress" class="block text-gray-700 mb-1">Address *</label>
                                    <input
                                            type="text"
                                            id="shippingAddress"
                                            bind:value={formData.shippingAddress}
                                            class="w-full p-2 border rounded-md {errors.shippingAddress && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                            class:error-field={errors.shippingAddress && formSubmitted}
                                    />
                                    {#if errors.shippingAddress && formSubmitted}
                                        <p class="text-red-500 text-sm mt-1">{errors.shippingAddress}</p>
                                    {/if}
                                </div>

                                <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                                    <div class="form-group">
                                        <label for="shippingCity" class="block text-gray-700 mb-1">City *</label>
                                        <input
                                                type="text"
                                                id="shippingCity"
                                                bind:value={formData.shippingCity}
                                                class="w-full p-2 border rounded-md {errors.shippingCity && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                                class:error-field={errors.shippingCity && formSubmitted}
                                        />
                                        {#if errors.shippingCity && formSubmitted}
                                            <p class="text-red-500 text-sm mt-1">{errors.shippingCity}</p>
                                        {/if}
                                    </div>

                                    <div class="form-group">
                                        <label for="shippingState" class="block text-gray-700 mb-1">State/Province *</label>
                                        <input
                                                type="text"
                                                id="shippingState"
                                                bind:value={formData.shippingState}
                                                class="w-full p-2 border rounded-md {errors.shippingState && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                                class:error-field={errors.shippingState && formSubmitted}
                                        />
                                        {#if errors.shippingState && formSubmitted}
                                            <p class="text-red-500 text-sm mt-1">{errors.shippingState}</p>
                                        {/if}
                                    </div>

                                    <div class="form-group">
                                        <label for="shippingZipCode" class="block text-gray-700 mb-1">ZIP/Postal Code *</label>
                                        <input
                                                type="text"
                                                id="shippingZipCode"
                                                bind:value={formData.shippingZipCode}
                                                class="w-full p-2 border rounded-md {errors.shippingZipCode && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                                class:error-field={errors.shippingZipCode && formSubmitted}
                                        />
                                        {#if errors.shippingZipCode && formSubmitted}
                                            <p class="text-red-500 text-sm mt-1">{errors.shippingZipCode}</p>
                                        {/if}
                                    </div>
                                </div>

                                <div class="form-group">
                                    <label for="shippingCountry" class="block text-gray-700 mb-1">Country *</label>
                                    <select
                                            id="shippingCountry"
                                            bind:value={formData.shippingCountry}
                                            class="w-full p-2 border rounded-md {errors.shippingCountry && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                            class:error-field={errors.shippingCountry && formSubmitted}
                                    >
                                        <option value="United States">United States</option>
                                        <option value="Canada">Canada</option>
                                        <option value="United Kingdom">United Kingdom</option>
                                        <option value="Australia">Australia</option>
                                        <option value="Germany">Germany</option>
                                        <option value="France">France</option>
                                    </select>
                                    {#if errors.shippingCountry && formSubmitted}
                                        <p class="text-red-500 text-sm mt-1">{errors.shippingCountry}</p>
                                    {/if}
                                </div>
                            </div>
                        {/if}
                    </div>

                    <!-- Payment Information -->
                    <div class="mb-8">
                        <h2 class="text-xl font-bold mb-4">Payment Method</h2>

                        <div class="mb-4">
                            <div class="flex items-center mb-2">
                                <input
                                        type="radio"
                                        id="creditCard"
                                        value="credit"
                                        bind:group={formData.paymentMethod}
                                        class="mr-2"
                                />
                                <label for="creditCard" class="text-gray-700">Credit Card</label>
                            </div>

                            <div class="flex items-center">
                                <input
                                        type="radio"
                                        id="paypal"
                                        value="paypal"
                                        bind:group={formData.paymentMethod}
                                        class="mr-2"
                                />
                                <label for="paypal" class="text-gray-700">PayPal</label>
                            </div>
                        </div>

                        {#if formData.paymentMethod === 'credit'}
                            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                                <div class="form-group md:col-span-2">
                                    <label for="cardNumber" class="block text-gray-700 mb-1">Card Number *</label>
                                    <input
                                            type="text"
                                            id="cardNumber"
                                            bind:value={formData.cardNumber}
                                            placeholder="1234 5678 9012 3456"
                                            class="w-full p-2 border rounded-md {errors.cardNumber && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                            class:error-field={errors.cardNumber && formSubmitted}
                                    />
                                    {#if errors.cardNumber && formSubmitted}
                                        <p class="text-red-500 text-sm mt-1">{errors.cardNumber}</p>
                                    {/if}
                                </div>

                                <div class="form-group md:col-span-2">
                                    <label for="cardName" class="block text-gray-700 mb-1">Name on Card *</label>
                                    <input
                                            type="text"
                                            id="cardName"
                                            bind:value={formData.cardName}
                                            class="w-full p-2 border rounded-md {errors.cardName && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                            class:error-field={errors.cardName && formSubmitted}
                                    />
                                    {#if errors.cardName && formSubmitted}
                                        <p class="text-red-500 text-sm mt-1">{errors.cardName}</p>
                                    {/if}
                                </div>

                                <div class="form-group">
                                    <label for="expiryDate" class="block text-gray-700 mb-1">Expiry Date *</label>
                                    <input
                                            type="text"
                                            id="expiryDate"
                                            bind:value={formData.expiryDate}
                                            placeholder="MM/YY"
                                            class="w-full p-2 border rounded-md {errors.expiryDate && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                            class:error-field={errors.expiryDate && formSubmitted}
                                    />
                                    {#if errors.expiryDate && formSubmitted}
                                        <p class="text-red-500 text-sm mt-1">{errors.expiryDate}</p>
                                    {/if}
                                </div>

                                <div class="form-group">
                                    <label for="cvv" class="block text-gray-700 mb-1">CVV *</label>
                                    <input
                                            type="text"
                                            id="cvv"
                                            bind:value={formData.cvv}
                                            placeholder="123"
                                            class="w-full p-2 border rounded-md {errors.cvv && formSubmitted ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-emerald-300"
                                            class:error-field={errors.cvv && formSubmitted}
                                    />
                                    {#if errors.cvv && formSubmitted}
                                        <p class="text-red-500 text-sm mt-1">{errors.cvv}</p>
                                    {/if}
                                </div>
                            </div>
                        {:else if formData.paymentMethod === 'paypal'}
                            <p class="text-gray-600 mb-4">You will be redirected to PayPal to complete your payment after reviewing your order.</p>
                        {/if}
                    </div>

                    <!-- Submit Button -->
                    <div>
                        <button
                                type="submit"
                                class="w-full bg-emerald-500 hover:bg-emerald-600 text-white py-3 rounded-md transition-colors"
                        >
                            Place Order
                        </button>
                    </div>
                </form>
            </div>

            <!-- Order Summary -->
            <div class="lg:col-span-1">
                <div class="bg-white rounded-lg shadow-sm p-6 sticky top-8">
                    <h2 class="text-xl font-bold mb-4">Order Summary</h2>

                    <!-- Cart Items -->
                    <div class="mb-6">
                        {#each cartItems as item}
                            <div class="flex items-center py-3 border-b border-gray-100 last:border-0">
                                <div class="w-16 h-16 bg-gray-50 rounded overflow-hidden mr-4">
                                    <img
                                            src={item.image}
                                            alt={item.name}
                                            class="w-full h-full object-cover"
                                    />
                                </div>
                                <div class="flex-grow">
                                    <h3 class="text-gray-800 font-medium">{item.name}</h3>
                                    <div class="flex justify-between text-sm text-gray-600">
                                        <span>Qty: {item.quantity}</span>
                                        <span>${(item.price * item.quantity).toFixed(2)}</span>
                                    </div>
                                </div>
                            </div>
                        {/each}
                    </div>

                    <!-- Totals -->
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

                    <!-- Return to Cart -->
                    <a
                            href="/cart"
                            class="text-emerald-600 hover:text-emerald-700 flex items-center text-sm"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                        </svg>
                        Return to Cart
                    </a>
                </div>
            </div>
        </div>
    </div>
</div>
