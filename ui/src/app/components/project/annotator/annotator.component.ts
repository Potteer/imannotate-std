import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ProjectService } from '@app/services/project.service';
import { Project } from '@app/classes/project';
import { Annotator } from '@app/classes/annotator';
import { BoundingBox } from '@app/classes/boundingbox';
import { Annotation } from '@app/classes/annotation';
import { ImageResult } from "@app/classes/imageresult";

import {NgbModal, ModalDismissReasons} from '@ng-bootstrap/ng-bootstrap';

@Component({
  selector: 'app-annotator',
  templateUrl: './annotator.component.html',
  styleUrls: ['./annotator.component.css']
})
export class AnnotatorComponent implements OnInit {
  annotator: Annotator;
  currentBox: BoundingBox;
  boxes = new Array<BoundingBox>();
  project = new Project();
  image: ImageResult ;
  idnumber: number;

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
    if ((this.currentBox.width - this.currentBox.x) * this.annotator.canvas.clientWidth < 0.1 ||
        (this.currentBox.height - this.currentBox.y) * this.annotator.canvas.clientHeight < 0.1
    ) {
      alert("Box too small !")
      return
    }

    this.currentBox.x = 0;
    this.currentBox.y= 0;
    this.currentBox.width = 0;
    this.currentBox.height = 0;
    this.currentBox.color = 'red';
    this.currentBox.id = 'label-' + Math.floor(Math.random() * 99999999).toString();;
    this.currentBox.old_color = null;
    this.currentBox.textColor = 'white';
    this.currentBox.hasAlpha = false;

    this.currentBox.label = label;
    this.boxes.push(this.currentBox);
    console.log("boxes após setLabel: ",this.boxes)
    this.annotator.setBoundingBoxes(this.boxes);
    console.log("annotator após setLabel: ",this.annotator)
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
