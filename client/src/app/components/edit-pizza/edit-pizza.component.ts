import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Ingredient } from '../../definitions/ingredient.class';
import { Category } from '../../definitions/category.class';
import { PizzaCanvasService } from '../../services/pizza-canvas.service';
import { PizzaDataService } from '../../services/pizza-data.service';
import { FormControl, FormArray, Validators, NgForm } from '@angular/forms';
import { Pizza } from '../../classes/pizza.class';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs';

@Component({
  selector: 'edit-create-pizza',
  templateUrl: './edit-pizza.component.html',
  styleUrls: ['./edit-pizza.component.sass'],
  interpolation: ["{$", "$}"]
})
export class EditPizzaComponent implements OnInit, OnDestroy {
  ingredients: Array<Ingredient>
  categories: Array<Category>
  categories_keys:Array<string>
  ingredients_keys:Array<string>
  connections:Subscription[]
  checkboxesFormArray: FormArray
  @ViewChild('editForm') public editForm: NgForm

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
    this.editForm.form.addControl("checkboxes", this.checkboxesFormArray)
  }

  changeCheckbox(ingredient: Ingredient) {
    var index = this.checkboxesFormArray.value.indexOf(ingredient)
    if (index !== -1) {
      this.checkboxesFormArray.removeAt(index)
    } else {
      this.checkboxesFormArray.push(new FormControl(ingredient))
    }
  }

  updatePizza() {
    this.connections.push(
      this.pizzaData.updatePizza(
          this.pizzaCanvas.getSavingData(),
          this.pizzaData.currentPizza.id
        ).subscribe(data => {
        this.router.navigate(["/"])
      })
    )
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
    hashCategories[this.pizzaData.currentPizza.category.name] = this.pizzaData.currentPizza.category
    this.pizzaCanvas.init(new Pizza({
      name: this.pizzaData.currentPizza.name,
      description: this.pizzaData.currentPizza.description,
      price: this.pizzaData.currentPizza.price,
      size: this.pizzaData.currentPizza.size,
      ingredients: this.pizzaData.currentPizza.ingredients,
      category: this.pizzaData.currentPizza.category
    }), hashIngredients, hashCategories)
    .then(() => this.addCheckboxes(this.ingredients
      .filter(ingredient => ~this.pizzaData.currentPizza.ingredients
        .findIndex(ingr=> ingredient.id == ingr.id)
        )
    ))
  }
}

