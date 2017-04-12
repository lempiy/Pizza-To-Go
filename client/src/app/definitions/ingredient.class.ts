import * as PIXI from 'pixi.js';
export class Ingredient {
  id?: number
  name: string
  description?: string
  price: number
  img_url: string
  checked?: boolean
  container?: PIXI.Container
}
