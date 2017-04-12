import { Point } from '../definitions/point.class';

export class Polygon {
  radius: number
  corners: number
  center: Point
  coords: Array<Point>
  constructor(center: Point, radius: number, corners: number) {
    this.radius = radius
    this.corners = corners
    this.center = center
    this.create()
  }
  private create() {
    if (this.corners < 2) {
      let shifttedPoint = {
        x: this.center.x + this.getRandomInt(-32, 32),
        y: this.center.y + this.getRandomInt(-32, 32)
      }
      this.coords = [shifttedPoint]
      return
    }

    this.coords = [];

    for (let i = 0; i <= this.corners - 1; i++) {
        this.coords.push(
          {
            x: Math.floor(this.center.x + this.radius * Math.cos(i * 2 * Math.PI / this.corners)),
            y: Math.floor(this.center.y + this.radius * Math.sin(i * 2 * Math.PI / this.corners))
          }
        )
    }
  }
  private getRandomInt(min:number, max:number) {
      min = Math.ceil(min);
      max = Math.floor(max);
      return Math.floor(Math.random() * (max - min + 1)) + min;
  }
}
