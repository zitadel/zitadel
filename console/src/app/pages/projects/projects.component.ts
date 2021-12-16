import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';

@Component({
  selector: 'cnsl-projects',
  templateUrl: './projects.component.html',
  styleUrls: ['./projects.component.scss'],
})
export class ProjectsComponent {
  public zitadelProjectId: string = '';
  public projectType$: BehaviorSubject<any> = new BehaviorSubject(ProjectType.PROJECTTYPE_OWNED);
  public ProjectType: any = ProjectType;
  public grid: boolean = true;
  constructor(private router: Router, mgmtService: ManagementService, breadcrumbService: BreadcrumbService) {
    mgmtService.getIAM().then((iam) => {
      this.zitadelProjectId = iam.iamProjectId;
    });

    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);
  }

  public addProject(): void {
    this.router.navigate(['/projects', 'create']);
  }

  public setType(type: ProjectType) {
    this.projectType$.next(type);
  }
}
