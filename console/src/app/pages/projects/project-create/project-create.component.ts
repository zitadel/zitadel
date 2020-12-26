import { Location } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Project, ProjectCreateRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
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
        private mgmtService: ManagementService,
        private _location: Location,
    ) { }

    public createSteps: number = 1;
    public currentCreateStep: number = 1;
    public ngOnInit(): void { }

    public saveProject(): void {
        this.mgmtService
            .CreateProject(this.project)
            .then((data: Project) => {
                this.router.navigate(['projects', data.getId()]);
            })
            .catch(error => {
                this.toast.showError(error);
            });
    }

    public close(): void {
        this._location.back();
    }
}
