import {type Writable, writable} from 'svelte/store';

export interface IcartItem {
    prodID: string;
    quantity: number;
}

function savedCart (): IcartItem[] {
    let cart: string | null = null;
    try {
        cart = localStorage.getItem("cart");
    } catch {
        console.log("Cart not found")
    }
    return cart? JSON.parse(cart) : [];
}

export const shoppingCart:Writable<IcartItem[]> = writable(savedCart());


export function addCartItem(newItem: IcartItem): void {
    shoppingCart.update(current => {
        const newCart:IcartItem[] = [ ...current, newItem ]

        const groupObj = newCart.reduce((r,{prodID,quantity}) =>
            (r[prodID] = (r[prodID]||0) + quantity, r), {});
        const grouped = Object.keys(groupObj).map(key => ({product:key, quantity: groupObj[key]}))


        console.log(grouped);
        localStorage.setItem("cart", JSON.stringify(grouped));
        return newCart;
    });
}