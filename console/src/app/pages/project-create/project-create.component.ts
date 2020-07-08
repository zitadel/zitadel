import { Location } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Project, ProjectCreateRequest } from 'src/app/proto/generated/management_pb';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-project-create',
    templateUrl: './project-create.component.html',
    styleUrls: ['./project-create.component.scss'],
})
export class ProjectCreateComponent implements OnInit {
    public project: ProjectCreateRequest.AsObject = new ProjectCreateRequest().toObject();

    constructor(
        private router: Router,
        private toast: ToastService,
        private projectService: ProjectService,
        private _location: Location,
    ) { }

    public createSteps: number = 1;
    public currentCreateStep: number = 1;
    public ngOnInit(): void { }

    public saveProject(): void {
        this.projectService
            .CreateProject(this.project)
            .then((data: Project) => {
                this.router.navigate(['projects', data.getId()]);
            })
            .catch(data => {
                this.toast.showError(data.message);
            });
    }

    public close(): void {
        this._location.back();
    }
}
