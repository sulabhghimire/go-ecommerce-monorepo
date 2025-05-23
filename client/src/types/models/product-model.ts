export interface CartModel {
  id: number;
  user_id: number;
  product_id: number;
  name: string;
  image_url: string;
  seller_id: number;
  price: number;
  qty: number;
  created_at: string;
  updated_at: string;
}

export interface ProductModel {
  id: number;
  name: string;
  description: string;
  category_id: string;
  image_url: string;
  price: number;
  stock: number;
  availability: boolean;
}
