import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpModule, Http, RequestOptions } from '@angular/http';

import { routing, appRoutesProviders } from './app.routing';
import { AppComponent } from './app.component';
import { NavbarComponent } from './components/navbar/navbar.component';
import { HomeComponent } from './components/home/home.component';
import { AuthComponent } from './components/auth/auth.component';
import { CreatePizzaComponent } from './components/create-pizza/create-pizza.component';
import { EditPizzaComponent } from './components/edit-pizza/edit-pizza.component';

import { AuthService } from './services/auth.service';
import { PizzaCanvasService } from './services/pizza-canvas.service';
import { PizzaDataService } from './services/pizza-data.service';
import { AuthHttp, AuthConfig } from 'angular2-jwt';
import { AuthGuard } from './auth.guard';
import { PizzaGuard } from './pizza.guard';
import { PizzaDetailsComponent } from './components/pizza-details/pizza-details.component';


export function authHttpServiceFactory(http: Http, options: RequestOptions) {
  return new AuthHttp(new AuthConfig({
    tokenName: 'token',
		tokenGetter: (() => localStorage.getItem('token')),
		globalHeaders: [{'Content-Type':'application/json'}],
    noTokenScheme: true
	}), http, options);
}

@NgModule({
  declarations: [
    AppComponent,
    NavbarComponent,
    HomeComponent,
    AuthComponent,
    CreatePizzaComponent,
    EditPizzaComponent,
    PizzaDetailsComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    routing,
    ReactiveFormsModule
  ],
  providers: [
    appRoutesProviders,
    AuthService,
    {
      provide: AuthHttp,
      useFactory: authHttpServiceFactory,
      deps: [Http, RequestOptions]
    },
    AuthGuard,
    PizzaGuard,
    PizzaCanvasService,
    PizzaDataService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
