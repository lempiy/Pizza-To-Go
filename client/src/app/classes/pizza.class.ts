import { Ingredient } from '../definitions/ingredient.class';
import { Category } from '../definitions/category.class';
export class Pizza {
  id?: number
  name: string
  description: string
  price: number
  image_url?: string
  ingredients: Array<Ingredient>
  size: number
  category: Category
  constructor(data) {
    this.id = data.id ? data.id : null
    this.name = data.name
    this.description = data.description
    this.price = data.price
    this.size = data.size
    this.image_url = data.image_url ? data.image_url : null
    this.ingredients = data.ingredients
    this.category = data.category
  }
}
