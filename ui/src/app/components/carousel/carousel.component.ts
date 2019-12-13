import { Component, ViewChild } from '@angular/core';
import { NgbCarousel, NgbSlideEvent, NgbSlideEventSource } from '@ng-bootstrap/ng-bootstrap';


@Component(
    {selector: 'ngbd-carousel-pause',
    templateUrl: './carousel.component.html'})
export class NgbdCarouselPause {
  images = [62, 83, 466, 965, 982, 1043, 738].map((n) => `https://picsum.photos/id/${n}/900/500`);

  paused = false;
  unpauseOnArrow = false;
  pauseOnIndicator = false;
  pauseOnHover = true;

  @ViewChild('carousel') carousel: NgbCarousel;

  togglePaused() {
    if (this.paused) {
      this.carousel.cycle();
    } else {
      this.carousel.pause();
    }
    this.paused = !this.paused;
  }

  onSlide(slideEvent: NgbSlideEvent) {
    if (this.unpauseOnArrow && slideEvent.paused &&
      (slideEvent.source === NgbSlideEventSource.ARROW_LEFT || slideEvent.source === NgbSlideEventSource.ARROW_RIGHT)) {
      this.togglePaused();
    }
    if (this.pauseOnIndicator && !slideEvent.paused && slideEvent.source === NgbSlideEventSource.INDICATOR) {
      this.togglePaused();
    }
  }
}



// foto1= 'http://127.0.0.1:9000/liveness/cpf1/mulher1.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=accessKey1%2F20191212%2F%2Fs3%2Faws4_request&X-Amz-Date=20191212T185526Z&X-Amz-Expires=432000&X-Amz-SignedHeaders=host&X-Amz-Signature=e73195d6f29e13878ce961a327f53db3212ab2ba0afb295c1011c1d3fbb93ee0';
// foto2= 'http://127.0.0.1:9000/liveness/cpf1/mulher2.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=accessKey1%2F20191212%2F%2Fs3%2Faws4_request&X-Amz-Date=20191212T185506Z&X-Amz-Expires=432000&X-Amz-SignedHeaders=host&X-Amz-Signature=f655b31fc51fc5ce9d5299072406832d43a18949fa7f3ba0d0a6c8c018dff367';
// foto3= 'http://127.0.0.1:9000/liveness/cpf1/mulher3.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=accessKey1%2F20191212%2F%2Fs3%2Faws4_request&X-Amz-Date=20191212T185542Z&X-Amz-Expires=432000&X-Amz-SignedHeaders=host&X-Amz-Signature=53a3c6aa63b51835a00ee1ce5d9ce2a3d053d286855fe2e674003b2a16b220dd';