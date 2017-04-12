import { Component, OnInit, OnDestroy } from '@angular/core';
import { AuthService } from "../../services/auth.service";
import { Subscription } from "rxjs";
import { Router } from '@angular/router';

@Component({
  selector: 'pizza-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.sass'],
  interpolation: ["{$", "$}"]
})
export class NavbarComponent implements OnInit, OnDestroy {
  conections: Subscription[]
  constructor(private authService:AuthService, private router: Router) {
    this.conections = [];
  }

  logout() {
    this.conections.push(
      this.authService.logout().subscribe(data => this.router.navigate(["/"]))
    )
  }

  ngOnInit() {
  }
  ngOnDestroy() {
    this.conections.forEach(subscr => subscr.unsubscribe())
  }

}
