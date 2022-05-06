import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Subject, takeUntil } from 'rxjs';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';

export enum GrantType {
  ORG = 'org',
  PROJECT = 'project',
}

@Component({
  selector: 'cnsl-grants',
  templateUrl: './grants.component.html',
  styleUrls: ['./grants.component.scss'],
})
export class GrantsComponent implements OnDestroy {
  public grantContext: UserGrantContext = UserGrantContext.NONE;
  public projectId: string = '';
  public UserGrantContext: any = UserGrantContext;
  public isZitadel: boolean = false;
  public destroy$: Subject<void> = new Subject();

  constructor(
    public translate: TranslateService,
    activatedRoute: ActivatedRoute,
    private mgmtService: ManagementService,
    private breadcrumbService: BreadcrumbService,
  ) {
    activatedRoute.data.pipe(takeUntil(this.destroy$)).subscribe((params) => {
      const { context } = params;
      this.grantContext = context;
      if (context === UserGrantContext.OWNED_PROJECT) {
        const projectId = activatedRoute.snapshot.paramMap.get('projectid');
        if (projectId) {
          this.projectId = projectId;

          this.mgmtService.getIAM().then((iam) => {
            this.isZitadel = iam.iamProjectId === this.projectId;

            const breadcrumbs = [
              new Breadcrumb({
                type: BreadcrumbType.ORG,
                routerLink: ['/org'],
              }),
              new Breadcrumb({
                type: BreadcrumbType.PROJECT,
                name: '',
                param: { key: 'projectid', value: this.projectId },
                routerLink: ['/projects', this.projectId],
                isZitadel: this.isZitadel,
              }),
            ];
            this.breadcrumbService.setBreadcrumb(breadcrumbs);
          });
        }
      } else if (context === UserGrantContext.NONE) {
        const bread: Breadcrumb = {
          type: BreadcrumbType.ORG,
          routerLink: ['/org'],
        };
        breadcrumbService.setBreadcrumb([bread]);
      }
    });
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
