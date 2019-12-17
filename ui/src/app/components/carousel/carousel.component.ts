import { Component, ViewChild } from '@angular/core';
import { NgbCarousel, NgbSlideEvent, NgbSlideEventSource, NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { Annotator } from '@app/classes/annotator';
import { Annotation } from '@app/classes/annotation';
import { ImageResult } from '@app/classes/imageresult';
import { Project } from '@app/classes/project';
import { BoundingBox } from '@app/classes/boundingbox';
import { ActivatedRoute } from '@angular/router';
import { ProjectService } from '@app/services/project.service';


@Component(
    {selector: 'ngbd-carousel-pause',
    templateUrl: './carousel.component.html'})
export class NgbdCarouselPause {
  //images = [62, 83, 466, 965, 982, 1043, 738].map((n) => `https://picsum.photos/id/${n}/900/500`);
  images = [];

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

export class AnnotatorComponent {
  annotator: Annotator;
  currentBox: BoundingBox;
  boxes = new Array<BoundingBox>();
  project = new Project();
  image: ImageResult;

  constructor(
    private route: ActivatedRoute,
    private projectService: ProjectService,
    private modalService: NgbModal,
  ) { }

  ngOnInit() {
    this.route.params.subscribe(param => {
      console.log('Param => ', param);
      this.projectService.getProject(param.name).subscribe(project => {
        this.project = project;
        this.annotator = new Annotator('#annotator');

        this.annotator.boxes.subscribe(box => {
          this.currentBox = box;
        });

        this.annotator.removeBox.subscribe(box => {
          const index: number = this.boxes.indexOf(box);
          if (index !== -1) {
            this.boxes.splice(index, 1);
          }
          this.annotator.setBoundingBoxes(this.boxes);
        });

        this.nextImage();
      });
    });
  }

  setLabel(label: string) {
    if (!this.currentBox) { return; }
    if ((this.currentBox.width - this.currentBox.x) * this.annotator.canvas.clientWidth < 20 ||
        (this.currentBox.height - this.currentBox.y) * this.annotator.canvas.clientHeight < 20
    ) {
      alert("Box too small !")
      return
    }

    this.currentBox.label = label;
    this.boxes.push(this.currentBox);
    this.annotator.setBoundingBoxes(this.boxes);
    this.currentBox = null;
  }


  saveAnnotation() {
    console.log(this.boxes);
    const annotation = new Annotation();
    annotation.image = this.image.name;
    annotation.boxes = this.boxes;
    console.log(annotation);
    this.projectService.saveAnnotation(this.project, annotation).subscribe(ann => {
      console.log("saved:", ann);
      this.nextImage();
    });
  }

  saveEmptyAnnotation(content) {
    this.modalService.open(content, {}).result.then(
      result => {
        console.log("result", result);
      },
      reason => {
        console.log("reason", reason);
      }
    );
  }

  nextImage() {
    // TODO: send box to server before to get next image
    this.projectService.getNextImage().subscribe(image => {
      console.log(image)
      this.image = image;
      this.boxes = new Array<BoundingBox>();
      this.annotator.loadImage(image.url);
    });
  }
}




// foto1= 'http://127.0.0.1:9000/liveness/cpf1/mulher1.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=accessKey1%2F20191212%2F%2Fs3%2Faws4_request&X-Amz-Date=20191212T185526Z&X-Amz-Expires=432000&X-Amz-SignedHeaders=host&X-Amz-Signature=e73195d6f29e13878ce961a327f53db3212ab2ba0afb295c1011c1d3fbb93ee0';
// foto2= 'http://127.0.0.1:9000/liveness/cpf1/mulher2.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=accessKey1%2F20191212%2F%2Fs3%2Faws4_request&X-Amz-Date=20191212T185506Z&X-Amz-Expires=432000&X-Amz-SignedHeaders=host&X-Amz-Signature=f655b31fc51fc5ce9d5299072406832d43a18949fa7f3ba0d0a6c8c018dff367';
// foto3= 'http://127.0.0.1:9000/liveness/cpf1/mulher3.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=accessKey1%2F20191212%2F%2Fs3%2Faws4_request&X-Amz-Date=20191212T185542Z&X-Amz-Expires=432000&X-Amz-SignedHeaders=host&X-Amz-Signature=53a3c6aa63b51835a00ee1ce5d9ce2a3d053d286855fe2e674003b2a16b220dd';