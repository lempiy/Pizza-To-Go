import { Injectable } from '@angular/core';
import { Pizza } from '../classes/pizza.class';
import { Ingredient } from '../definitions/ingredient.class';
import { Category } from '../definitions/category.class';
import { Point } from '../definitions/point.class';
import { PizzaData } from '../definitions/pizza-data.class';
import { FillingField } from '../classes/filling-field.class';
import * as PIXI from 'pixi.js';

@Injectable()
export class PizzaCanvasService {
  pizza: Pizza
  angles: Array<number>
  ingredients: {[name: string]: Ingredient}
  categories: {[name: string]: Category}

  private rotationMap:Array<number>
  private renderer: PIXI.WebGLRenderer | PIXI.CanvasRenderer
  private rootContainer: PIXI.Container
  private ingrContainer: PIXI.Container
  private loader:any;
  private res:PIXI.loaders.Loader;
  private ready:boolean
  private animFrameID: number
  //private ingredients_sprites: {[name: string]: PIXI.Sprite}

  private basicSprite: PIXI.Sprite

  constructor() {

  }
  init(pizza: Pizza, ingredients: {[name: string]: Ingredient}, categories: {[name: string]: Category}):Promise<any> {
    this.rootContainer = new PIXI.Container();
    this.rootContainer.position.set(0, 0);
    this.rootContainer.width = 400;
    this.rootContainer.height = 400;

    this.pizza = pizza
    this.ingredients = ingredients
    this.categories = categories
    this.createRotationMap()

    this.renderer = new PIXI.WebGLRenderer(400, 400, { antialias: true, transparent: true, preserveDrawingBuffer: true });
    document.querySelector(".render-container").appendChild(this.renderer.view);

    if (this.res && this.loader) {
        return Promise.resolve(this.run(this.res, this.loader));
    }
    return new Promise((resolve, reject) => {
        let ld = PIXI.loader.add("pizza", '/assets/images/pizza.png')
        Object.keys(this.ingredients).forEach(ingredient => {
          ld.add(this.ingredients[ingredient].name, this.ingredients[ingredient].img_url)
        })
        ld.load((res, loader) => {
          this.run(res, loader)
          resolve()
        })
    })
  }
  run(res:PIXI.loaders.Loader, loader:any) {
    this.res = res
    this.loader = loader
    this.drawBase()
    this.initIngredients()
    this.drawIngredients()
    this.changeSize(this.pizza.size)
    this.ready = true
    this.animate()
  }

  private animate() {
      let self = this;

      self.renderer.render(this.rootContainer)
      let prev = Date.now()
      self.animFrameID = requestAnimationFrame(function(){
          let now = Date.now(),
              delta = now / prev
          self.spin(delta)
          self.animate.call(self)
      });
  }

  stop() {
    cancelAnimationFrame(this.animFrameID)
  }

  private drawBase() {
    this.basicSprite = PIXI.Sprite.fromImage("pizza")
    this.basicSprite.anchor.set(0.5, 0.5)
    this.basicSprite.position.set(
      Math.floor(400 * 0.5),
      Math.floor(400 * 0.5)
    )
    this.rootContainer.addChildAt(this.basicSprite, 0)
  }

  private initIngredients() {
    this.ingrContainer = new PIXI.Container()
    this.ingrContainer.width = 400
    this.ingrContainer.height = 400
    this.rootContainer.addChildAt(this.ingrContainer, 1)
    this.setRootPivot()
  }

  private drawIngredients() {
    this.pizza.ingredients.forEach(ingrd => {
      this.ingredients[ingrd.name].checked = true
      this.addIngredient(ingrd.name)
    })
  }

  addIngredient(name:string) {
    if(this.ingredients[name].container) {
      this.ingredients[name].container.visible = true
    } else {
      this.drawIngredientLayer(name, this.countDrawedIngredients())
    }
  }

  removeIngredient(name:string) {
    this.ingredients[name].container.visible = false
  }

  private drawIngredientLayer(name:string, order:number) {
    let fillmap = new FillingField(340, 30, {x: 200, y: 200});
    let fullCont = new PIXI.Container();
    let portion = (Math.PI * 2) / fillmap.map.length
    fillmap.map.forEach((poly, index) => {
      let ingrCont = new PIXI.Container();
      poly.coords.forEach(point => {
        let sprite = PIXI.Sprite.fromImage(name)
        sprite.anchor.set(0.5, 0.5)
        sprite.position.set(point.x, point.y)
        ingrCont.addChild(sprite)
      })
      ingrCont.position.set(200, 200)
      ingrCont.pivot.set(200, 200)
      fullCont.addChild(ingrCont)
      ingrCont.rotation = portion * (index + 1)
    })
    fullCont.position.set(200, 200)
    fullCont.pivot.set(200, 200)
    fullCont.rotation = this.rotationMap[order]
    this.ingrContainer.addChild(fullCont)
    this.ingredients[name].container = fullCont
  }

  private createRotationMap() {
    let length = Object.keys(this.ingredients).length
    this.rotationMap = [];
    for(let i = 1; i <= length; i++)
    {
      this.rotationMap.push(Math.PI * 2 / length * i)
    }
    this.rotationMap.sort((a, b) => Math.random() - 0.5)
  }

  private countDrawedIngredients() {
    return Object.keys(this.ingredients).filter(key => this.ingredients[key].container).length
  }

  private setRootPivot() {
    this.rootContainer.position.set(200, 200)
    this.rootContainer.pivot.set(200, 200)
  }

  changeSize(size:number) {
    switch(size) {
      case 30:
        this.rootContainer.scale.set(0.8, 0.8)
        break
      case 45:
        this.rootContainer.scale.set(0.9, 0.9)
        break
      case 60:
        this.rootContainer.scale.set(1, 1)
        break
    }
  }

  getSavingData():PizzaData {
    let pizzaIngredientsIds = [];
    Object.keys(this.ingredients).forEach(key => {
      if (this.ingredients[key].checked) {
        pizzaIngredientsIds.push(this.ingredients[key].id)
      }
    })
    return {
      name: this.pizza.name,
      size: this.pizza.size,
      description: this.pizza.description,
      category: this.pizza.category.id,
      ingredients: pizzaIngredientsIds,
      encodedImage: this.getSnapShot()
    }
  }

  private spin(delta: number) {
    this.rootContainer.rotation -= 0.003 * delta;
  }

  private getSnapShot():string {
    return this.renderer.view.toDataURL()
  }
}
