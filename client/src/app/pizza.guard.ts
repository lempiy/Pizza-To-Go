import { Injectable } from '@angular/core';
import { Router, ActivatedRouteSnapshot, RouterStateSnapshot, CanActivate } from '@angular/router';
import { PizzaDataService } from './services/pizza-data.service';
import { Observable } from 'rxjs';

@Injectable()

export class PizzaGuard implements CanActivate {
    constructor(private pizzaData:PizzaDataService, private router:Router){

    }
    canActivate(
      next: ActivatedRouteSnapshot,
      state: RouterStateSnapshot
      ): Observable<boolean>|Promise<boolean>|boolean {
        return this.pizzaData.getIsMyPizzaPizza(next.params['id']).map((data: boolean) => {
          if (!data) {
            this.router.navigate(["/"])
          }
          return data
        })
    }
}
