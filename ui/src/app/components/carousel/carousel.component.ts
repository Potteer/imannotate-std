import {Component, ViewChild} from '@angular/core';
import {NgbCarouselConfig} from '@ng-bootstrap/ng-bootstrap';

@Component({
  selector: 'carousel',
  templateUrl: './carousel.component.html',
  providers: [NgbCarouselConfig] // add NgbCarouselConfig to the component providers
})
export class NgbdCarouselConfig {
  
   @ViewChild('carousel') carousel: any;
  
   paused=false;
  
  constructor(config: NgbCarouselConfig) {
    // customize default values of carousels used by this component tree
    config.interval = 500;
    config.keyboard = false;
    config.wrap = true;
    config.pauseOnHover = false;
  }

  togglePaused() {
    if (this.paused) {
      this.carousel.cycle();
      this.carousel.interval = 500;
    } else {
      this.carousel.pause();
      this.carousel.interval = 0;
    }
    this.paused = !this.paused;
  }

}