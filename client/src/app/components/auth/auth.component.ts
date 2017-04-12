import { Component, OnInit } from '@angular/core';
import { Response } from '@angular/http';
import { Router, ActivatedRoute } from '@angular/router';
import { FormControl } from '@angular/forms';
import { Subscription } from "rxjs"
import { AuthService } from "../../services/auth.service";

class User {
  username:string;
  password:string;
  email:string;
}

@Component({
  selector: 'pizza-auth',
  templateUrl: './auth.component.html',
  styleUrls: ['./auth.component.sass'],
  interpolation: ["{$", "$}"]
})
export class AuthComponent implements OnInit {
  private sub: Subscription
  public state: string
  public showLogin: boolean
  public showSignup: boolean

  user: User = new User()

  constructor(
    private router: Router,
    private activedRoute: ActivatedRoute,
    private authService: AuthService) {
  }

  login() {
    this.authService.login(this.user.username, this.user.password).subscribe(data => {
      this.router.navigate(["/"])
    })
  }

  signup() {
    this.authService.signup(this.user.username, this.user.email, this.user.password).subscribe(data => {
      this.router.navigate(["/"])
    }, (error: Response) => {
      let data = error.json()
      this.authService.errorMessage = data.message
    })
  }

  static	emailValidator(control:	FormControl)	{
    return	/^\w+@\w+\.\w{2,4}$/g.test(control.value)	?	null:	{	notAllowed:	true	};
  }

  emailCtrl:	FormControl	=	new	FormControl('',	AuthComponent.emailValidator);

  onKeyUp(event: KeyboardEvent) {
    if (this.authService.errorMessage) {
      this.authService.errorMessage = null
    }
  }

  ngOnInit() {
    this.sub = this.activedRoute
        .queryParams
        .subscribe(params => {
            this.state = params['auth']
            if (this.state === 'login') {
              this.showSignup = false
              this.showLogin = true
            } else if(this.state === 'signup') {
              this.showLogin = false
              this.showSignup = true
            } else {
              this.showLogin = false
              this.showSignup = false
            }
    });
  }

}
