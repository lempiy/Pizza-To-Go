import { Category } from "./category.class";

class IngredientData {
  id: number
  name: string
  img_url: string
  price: number
}
class Author {
  id: number
  name: string
  created: string
}
export class PizzaOrder {
  id: number
  name: string
  author: Author
  category: Category
  size: number
  description: string
  created_date: string
  updated_date: string
  img_url: string
  price: number
  ingredients: Array<IngredientData>
}

