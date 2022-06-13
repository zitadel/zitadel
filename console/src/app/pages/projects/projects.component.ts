import { Component } from '@angular/core';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { BehaviorSubject, take } from 'rxjs';
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
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private activatedRoute: ActivatedRoute,
    public mgmtService: ManagementService,
    breadcrumbService: BreadcrumbService,
  ) {
    this.activatedRoute.queryParams.pipe(take(1)).subscribe((params: Params) => {
      const type = params.type;
      if (type && type === 'owned') {
        this.setType(ProjectType.PROJECTTYPE_OWNED);
      } else if (type && type === 'granted') {
        this.setType(ProjectType.PROJECTTYPE_GRANTED);
      }
    });
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
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: {
        type:
          type === ProjectType.PROJECTTYPE_OWNED ? 'owned' : type === ProjectType.PROJECTTYPE_GRANTED ? 'granted' : 'owned',
      },
      replaceUrl: true,
      queryParamsHandling: 'merge',
      skipLocationChange: false,
    });
  }
}
