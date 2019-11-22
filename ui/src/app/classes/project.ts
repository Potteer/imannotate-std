export class Project {
  server: string;
  id: string;
  name: string;
  banner: string;
  owner: string;
  description: string;
  tags: Array<String>;
  imageProvider: string;
  imageProviderOptions: any;

  constructor() {
    this.server = "";
    this.id = "";
    this.name = "";
    this.owner = "";
    this.description = "";
    this.tags = new Array<string>();
    this.imageProvider = "";
    this.imageProviderOptions = {};
  }
}

export class ProjectStat {
  project: Project
  annotations: number
  contributors: number
}
