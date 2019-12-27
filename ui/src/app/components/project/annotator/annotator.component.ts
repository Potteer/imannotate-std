import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { ProjectService } from '@app/services/project.service';
import { Project } from '@app/classes/project';
import { Annotator } from '@app/classes/annotator';
import { Annotation } from '@app/classes/annotation';
import { ImageResult } from "@app/classes/imageresult";
import { UserService } from '@app/services/user.service';
import { User } from '@app/classes/user';


@Component({
  selector: 'app-annotator',
  templateUrl: './annotator.component.html',
  styleUrls: ['./annotator.component.css'],
})

export class AnnotatorComponent implements OnInit {
  annotator: Annotator;
  annotation: Annotation;
  project = new Project();
  image: ImageResult;
  label;
  boxes = new Array();
  public user: User;


  constructor(
    private route: ActivatedRoute,
    private projectService: ProjectService,
    private userservice: UserService,
  ) { this.user = new User();}

  ngOnInit() {
    this.route.params.subscribe(param => {
      console.log('Param => ', param);
      this.projectService.getProject(param.name).subscribe(project => {
        this.project = project;
        this.annotator = new Annotator('#annotator');

        this.annotator.boxes.subscribe(box => {
          this.boxes = box;
        });
        this.nextImage();
      });
    });
  }

  setLabel(label: string) {
    label = label;
    console.log("Antes: ", label)
  }

  saveAnnotation() {
    const annotation = new Annotation();
    annotation.image = this.image.name;
    annotation.username = this.userservice.currentUser.username;
    annotation.label = this.label;
    annotation.timestamp = Math.round(new Date().getTime()/1000) - Math.round(new Date().getTimezoneOffset() * 60 );
    this.projectService.saveAnnotation(this.project, annotation).subscribe(ann => {
      console.log("saved:", ann);
      console.log("label:", annotation.label);
      console.log("user:", annotation.username);
      console.log("Timestamp:", annotation.timestamp);
      this.nextImage();
    });
  }

  //JUntar tudo no boxes e mandar o boxes

  nextImage() {
    // TODO: send box to server before to get next image
    this.projectService.getNextImage().subscribe(image => {
      console.log("Next Image: ",image)
      this.image = image;
      //this.annotator.loadImage(image.url);
    });
  }
}


