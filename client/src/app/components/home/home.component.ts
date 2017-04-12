import { Component, OnInit, OnDestroy } from '@angular/core';
import { PizzaDataService } from '../../services/pizza-data.service';
import { AuthService } from '../../services/auth.service';
import { PizzaOrder } from '../../definitions/pizza-order.class';
import { Subscription } from 'rxjs';

@Component({
  selector: 'pizza-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.sass'],
  interpolation: ["{$", "$}"]
})
export class HomeComponent implements OnInit, OnDestroy {
  pizzas: Array<PizzaOrder>
  connections: Array<Subscription>
  selectedPizza: PizzaOrder
  showDetails: boolean
  constructor(private pizzaData: PizzaDataService, private auth: AuthService) {
    this.connections = []
    this.showDetails = false
  }

  isMyPizza(authorID:number) {
    return authorID === this.auth.userID
  }

  selectDetails(pizza:PizzaOrder) {
    this.selectedPizza = pizza
    this.showDetails = true
  }

  onClose(close: boolean) {
    this.showDetails = close
  }

  acceptPizza(id: number, index: number) {
    this.connections.push(
      this.pizzaData.acceptPizza(id).subscribe(data => {
        this.pizzaData.pizzas.splice(index, 1)
      })
    )
  }

  deletePizza(id: number, index: number) {
    this.connections.push(
      this.pizzaData.deletePizza(id).subscribe(data => {
        this.pizzaData.pizzas.splice(index, 1)
      })
    )
  }

  ngOnInit() {
    this.connections.push(
      this.subscribeToPizzas()
    )
  }

  ngOnDestroy() {
    this.connections.forEach(c => c.unsubscribe())
  }

  subscribeToPizzas() {
    return this.pizzaData.getPizzas().subscribe(pizzas => {
      this.pizzaData.pizzas = this.pizzaData.pizzas.map(pizza => {
        //TODO: MOVE THIS TRASH TO PIPE | USE OBSERVABLE DIRECTLY WITH async
        let created = new Date(pizza.created_date)
        pizza.created_date = `${created.getHours()}:${created.getMinutes()}:${
          created.getSeconds()} ${this.getFormattedDate(created.getDate(), created.getMonth() + 1)}`
        let updated = new Date(pizza.updated_date)
        pizza.updated_date = `${updated.getHours()}:${updated.getMinutes()}:${
          updated.getSeconds()} ${this.getFormattedDate(updated.getDate(), updated.getMonth() + 1)}`
        return pizza
      })
    })
  }
  private getFormattedDate(date: number, month: number) {
    let datestring = date > 9 ? "" + date : "0" + date
    let monthstring = month > 9 ? "" + month : "0" + month
    return `${datestring}.${monthstring}`
  }

}
