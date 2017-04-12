import { Component } from '@angular/core';

@Component({
  selector: 'pizza-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.sass'],
  interpolation: ["{$", "$}"]
})
export class AppComponent {
}
