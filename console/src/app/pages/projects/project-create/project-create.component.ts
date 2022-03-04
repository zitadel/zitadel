import { Location } from '@angular/common';
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AddProjectRequest, AddProjectResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-project-create',
  templateUrl: './project-create.component.html',
  styleUrls: ['./project-create.component.scss'],
})
export class ProjectCreateComponent {
  public project: AddProjectRequest.AsObject = new AddProjectRequest().toObject();

  constructor(
    private router: Router,
    private toast: ToastService,
    private mgmtService: ManagementService,
    private _location: Location,
    breadcrumbService: BreadcrumbService,
  ) {
    const iambread = new Breadcrumb({
      type: BreadcrumbType.IAM,
      name: 'IAM',
      routerLink: ['/system'],
    });
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([iambread, bread]);
  }

  public createSteps: number = 1;
  public currentCreateStep: number = 1;

  public saveProject(): void {
    this.mgmtService
      .addProject(this.project)
      .then((resp: AddProjectResponse.AsObject) => {
        this.router.navigate(['projects', resp.id]);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public close(): void {
    this._location.back();
  }
}
