export interface PetProfile {
    id: number;
    picture: string;
    name: string;
    date_of_birth?: string;
    tag_id?: string;
    city_license?: string;
    neutered?: boolean;
    vaccinated?: boolean;
    illnesses?: string[];
}

export interface PurchaseHistory {
    product_sku: string;
    registration_date: string; // ISO date format (YYYY-MM-DD)
    tag_id?: string; // Optional: Some products may not have a tag
};

export interface UserProfile {
    name: string;
    email: string;
    address: string;
    phone: string;
    joined_on: string;
    purchase_history: PurchaseHistory[];
    pets: PetProfile[];
}