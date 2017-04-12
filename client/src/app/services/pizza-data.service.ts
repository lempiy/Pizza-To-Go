import { Injectable } from '@angular/core';
import { Http } from '@angular/http'
import { AuthHttp } from 'angular2-jwt';
import { Observable } from 'rxjs';
import { PizzaOrder } from '../definitions/pizza-order.class';
import 'rxjs/add/operator/map';

@Injectable()
export class PizzaDataService {
  pendingRequest:boolean
  pizzas: Array<PizzaOrder>
  currentPizza: PizzaOrder
  constructor(private authHttp: AuthHttp, private http: Http) {

  }
  getCategories():Observable<any> {
    this.pendingRequest = true
    return this.authHttp.get('/api/categories/')
      .map(res => {
            this.pendingRequest = false;
            if (res.status < 400) {
              return res.json()
            }
        });
  }

  getPizzas():Observable<PizzaOrder[]> {
    this.pendingRequest = true
    return this.http.get('/api/get-pizzas/')
      .map(res => {
            this.pendingRequest = false;
            if (res.status < 400) {
              let pizzas = res.json()
              this.pizzas = pizzas
              return pizzas
            }
        });
  }

  getIngredients():Observable<any> {
    this.pendingRequest = true
    return this.authHttp.get('/api/ingredients/')
      .map(res => {
            this.pendingRequest = false;
            if (res.status < 400) {
              return res.json()
            }
        });
  }

  getCategoriesAndIngredients():Observable<any> {
    return Observable.forkJoin([this.getCategories(), this.getIngredients()]).map(res => {
      return {
        "categories": res[0],
        "ingredients": res[1]
      }
    })
  }

  postPizza(pizzaData: any):Observable<any> {
    this.pendingRequest = true
    return this.authHttp.post('/api/save-pizza/', pizzaData)
      .map(res => {
            this.pendingRequest = false;
            if (res.status < 400) {
              return res.json()
            }
        });
  }

  updatePizza(pizzaData: any, id: number):Observable<any> {
    this.pendingRequest = true
    return this.authHttp.put('/api/update-pizza/?id=' + id, pizzaData)
      .map(res => {
            this.pendingRequest = false;
            if (res.status < 400) {
              return res.json()
            }
        });
  }

  acceptPizza(id: number):Observable<any> {
    this.pendingRequest = true
    return this.authHttp.put('/api/accept-pizza/', {id: id})
      .map(res => {
            this.pendingRequest = false;
            if (res.status < 400) {
              return res.json()
            }
        });
  }

  deletePizza(id: number):Observable<any> {
    this.pendingRequest = true
    return this.authHttp.put('/api/delete-pizza/', {id: id})
      .map(res => {
            this.pendingRequest = false;
            if (res.status < 400) {
              return res.json()
            }
        });
  }

  getIsMyPizzaPizza(id: number):Observable<boolean> {
    this.pendingRequest = true
    return this.authHttp.get('/api/get-pizza-by-id/?id=' + id)
      .map(res => {
            this.pendingRequest = false;
            if (res.status < 400) {
              let data = res.json()
              this.currentPizza = data.pizza && data.pizza.id ? data.pizza : null
              return data.is_my_pizza
            }
        });
  }
}
