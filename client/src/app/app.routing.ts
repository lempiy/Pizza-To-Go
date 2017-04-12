import { ModuleWithProviders } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './components/home/home.component';
import { CreatePizzaComponent } from './components/create-pizza/create-pizza.component';
import { EditPizzaComponent } from './components/edit-pizza/edit-pizza.component';
import { AuthGuard } from './auth.guard';
import { PizzaGuard } from './pizza.guard';

const appRoutes:Routes = [
    {
        path: '',
        component: HomeComponent
    },
    {
        path: 'create',
        component: CreatePizzaComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'edit/:id',
        component: EditPizzaComponent,
        canActivate: [AuthGuard, PizzaGuard]
    }
]

export const appRoutesProviders: any[] = [];
export const routing: ModuleWithProviders = RouterModule.forRoot(appRoutes);
