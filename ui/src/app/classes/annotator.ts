import { Subject, Observable } from "rxjs";

export class Annotator {
  private image: HTMLImageElement;
  private element: Element;
  private boxeSubject;
  public boxes;

  constructor(element) {  
    this.boxeSubject = new Subject();
    this.boxes = this.boxeSubject.asObservable();

    if (typeof (element) === 'string') {
      this.element = document.querySelector(element);
    } else {
      this.element = element;
    }

  }

  toJson() {
    const boxes = [];
    (box) => {
      boxes.push({
        Label: box.label,
        Username: box.username,
        Timestamp: box.timestamp,
      });
    };
    return JSON.stringify({
      image: this.image.src,
      boxes: boxes,
    });
  }

  //Verificar como chamar 3x
  loadImage(url) {
    const image = new Image();
    image.addEventListener('load', () => {});
    image.src = url;
    this.image = image;
  }
}
