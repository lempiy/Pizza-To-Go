import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Ingredient } from '../../definitions/ingredient.class';
import { Category } from '../../definitions/category.class';
import { PizzaCanvasService } from '../../services/pizza-canvas.service';
import { PizzaDataService } from '../../services/pizza-data.service';
import { Pizza } from '../../classes/pizza.class';
import { NgForm, FormControl, FormArray, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs';

@Component({
  selector: 'pizza-create-pizza',
  templateUrl: './create-pizza.component.html',
  styleUrls: ['./create-pizza.component.sass'],
  interpolation: ["{$", "$}"]
})
export class CreatePizzaComponent implements OnInit, OnDestroy {
  ingredients: Array<Ingredient>
  categories: Array<Category>
  categories_keys:Array<string>
  ingredients_keys:Array<string>
  connections:Subscription[]
  checkboxesFormArray: FormArray
  @ViewChild('createForm') public createForm: NgForm

  constructor(private pizzaCanvas: PizzaCanvasService,
              private router: Router,
              private pizzaData: PizzaDataService) {

  }

  onSizeChange(value: number) {
    this.pizzaCanvas.pizza.price = value * this.pizzaCanvas.pizza.price / this.pizzaCanvas.pizza.size
    this.pizzaCanvas.pizza.size = value
    this.pizzaCanvas.changeSize(value)
  }

  onChangeIngredient(ingredientName:string, checked:boolean) {
    if (checked) {
      this.pizzaCanvas.addIngredient(ingredientName)
      this.pizzaCanvas.pizza.price += (this.pizzaCanvas.pizza.size / 60 * this.pizzaCanvas.ingredients[ingredientName].price)
    } else {
      this.pizzaCanvas.removeIngredient(ingredientName)
      this.pizzaCanvas.pizza.price -= (this.pizzaCanvas.pizza.size / 60 * this.pizzaCanvas.ingredients[ingredientName].price)
    }
    this.changeCheckbox(this.pizzaCanvas.ingredients[ingredientName])

  }

  private addCheckboxes(ingredients: Ingredient[]) {
    this.checkboxesFormArray = new FormArray(ingredients.map(ingredient => new FormControl(ingredient)),
      Validators.compose([Validators.required, Validators.minLength(3)]))
    this.createForm.form.addControl("checkboxes", this.checkboxesFormArray)
  }

  changeCheckbox(ingredient: Ingredient) {
    var index = this.checkboxesFormArray.value.indexOf(ingredient)
    if (index !== -1) {
      this.checkboxesFormArray.removeAt(index)
    } else {
      this.checkboxesFormArray.push(new FormControl(ingredient))
    }
  }

  submitPizza() {
    this.pizzaData.postPizza(this.pizzaCanvas.getSavingData()).subscribe(data => {
      this.router.navigate(["/"])
    })
  }

  ngOnDestroy() {
    this.pizzaCanvas.stop()
    this.connections.forEach(c => c.unsubscribe())
  }

  ngOnInit() {
    this.connections = []
    this.connections.push(
      this.pizzaData.getCategoriesAndIngredients().subscribe(data => {
        this.categories = data.categories
        this.ingredients = data.ingredients
        this.initData()
      })
    )
  }
  initData() {
    let hashCategories: {[name: string]: Category} = {}
    let hashIngredients: {[name: string]: Ingredient} = {}

    this.categories.forEach(category => {
      hashCategories[category.name] = category
    })

    this.ingredients.forEach(ingredient => {
      hashIngredients[ingredient.name] = ingredient
      hashIngredients[ingredient.name].checked = false
    })
    this.categories_keys = Object.keys(hashCategories)
    this.ingredients_keys = Object.keys(hashIngredients)
    this.pizzaCanvas.init(new Pizza({
      name: "",
      description: "",
      price: 10,
      size: 60,
      ingredients: [],
      category: this.categories.find((category) => category.is_default)
    }), hashIngredients, hashCategories)
    .then(() => this.addCheckboxes([]))
  }
}

