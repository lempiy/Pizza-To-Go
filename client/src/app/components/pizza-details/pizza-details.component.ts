import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { PizzaOrder } from '../../definitions/pizza-order.class';

@Component({
  selector: 'pizza-pizza-details',
  templateUrl: './pizza-details.component.html',
  styleUrls: ['./pizza-details.component.sass'],
  interpolation: ["{$", "$}"]
})
export class PizzaDetailsComponent implements OnInit {
  @Input() pizza: PizzaOrder
  @Output() onClose = new EventEmitter<boolean>();
  constructor() { }

  close() {
    this.onClose.emit(false)
  }

  ngOnInit() {
  }

}
