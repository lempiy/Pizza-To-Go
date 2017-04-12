import { Point } from '../definitions/point.class';
import { Polygon } from './polygon.class';

export class FillingField {
  areaSize: number
  itemSize: number
  center: Point
  map: Polygon[]
  constructor(areaSize: number, itemSize: number, center: Point) {
    this.areaSize = areaSize
    this.itemSize = itemSize
    this.center = center
    this.create()
  }
  private create() {
    let step = Math.floor(this.itemSize * 0.9)
    let amountOfSteps = Math.floor((this.areaSize * 0.5) / step)
    this.map = [];

    for (let i = 1; i < amountOfSteps; i++)
    {
      this.map.push(
        new Polygon(this.center, (this.areaSize * 0.5) - (step * i), amountOfSteps - i)
      )
    }
  }
}
