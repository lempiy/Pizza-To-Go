import { Injectable } from '@angular/core';
import { Router, ActivatedRouteSnapshot, RouterStateSnapshot, CanActivate } from '@angular/router';
import { AuthService } from './services/auth.service';

@Injectable()

export class AuthGuard implements CanActivate {
    constructor(private auth:AuthService, private router:Router){

    }
    canActivate(next: ActivatedRouteSnapshot, state: RouterStateSnapshot) {
        if(this.auth.loggedIn()) {
            return true;
        } else {
            console.log("User is not authorized")
            this.router.navigate(['/']);
            return false;
        }
    }
}
