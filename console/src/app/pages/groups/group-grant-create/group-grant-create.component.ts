import { Location } from '@angular/common';
import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Subject, takeUntil } from 'rxjs';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { GrantedProject, Project } from 'src/app/proto/generated/zitadel/project_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { Group } from 'src/app/proto/generated/zitadel/group_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { StorageKey, StorageLocation, StorageService } from 'src/app/services/storage.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-group-grant-create',
  templateUrl: './group-grant-create.component.html',
  styleUrls: ['./group-grant-create.component.scss'],
})
export class GroupGrantCreateComponent implements OnDestroy {
  public context!: UserGrantContext;

  public org?: Org.AsObject;
  public groupIds: string[] = [];

  public project?: Project.AsObject;
  public grantedProject?: GrantedProject.AsObject;

  public rolesList: string[] = [];

  public createSteps: number = 2;
  public currentCreateStep: number = 1;

  public UserGrantContext: any = UserGrantContext;

  public group?: Group.AsObject;

  public editState: boolean = false;
  private destroy$: Subject<void> = new Subject();

  constructor(
    private toast: ToastService,
    private _location: Location,
    private route: ActivatedRoute,
    private mgmtService: ManagementService,
    private storage: StorageService,
    breadcrumbService: BreadcrumbService,
  ) {
    breadcrumbService.setBreadcrumb([
      new Breadcrumb({
        type: BreadcrumbType.ORG,
        routerLink: ['/org'],
      }),
    ]);
    this.route.params.pipe(takeUntil(this.destroy$)).subscribe((params: Params) => {
      const { projectid, grantid, userid } = params;
      this.context = UserGrantContext.NONE;

      this.groupIds = userid ? [userid] : [];

      if (projectid && !grantid) {
        this.context = UserGrantContext.OWNED_PROJECT;

        this.mgmtService
          .getProjectByID(projectid)
          .then((resp) => {
            if (resp.project) {
              this.project = resp.project;
            }
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      } else if (projectid && grantid) {
        this.context = UserGrantContext.GRANTED_PROJECT;
        this.mgmtService
          .getGrantedProjectByID(projectid, grantid)
          .then((resp) => {
            if (resp.grantedProject) {
              this.grantedProject = resp.grantedProject;
            }
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      } else if (this.groupIds && this.groupIds.length === 1) {
        this.mgmtService
          .getGroupByID(this.groupIds[0])
          .then((resp) => {
            if (resp.group) {
              this.group = resp.group;
            }
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      }
    });

    const temporg = this.storage.getItem<Org.AsObject>(StorageKey.organization, StorageLocation.session);
    if (temporg) {
      this.org = temporg;
    }
  }

  public close(): void {
    this._location.back();
  }

  public addGrant(): void {
      let tempGrantId: string = '';
      if (this.grantedProject?.grantId) {
        tempGrantId = this.grantedProject.grantId;
      }
      const promn = this.groupIds.map((id) =>
        this.mgmtService.addGroupGrant(
          id,
          this.rolesList,
          this.project ? this.project.id : this.grantedProject ? this.grantedProject.projectId : '',
          tempGrantId,
        ),
      );
      Promise.all(promn)
        .then(() => {
          this.toast.showInfo('GRANTS.TOAST.UPDATED', true);
          this.close();
        })
        .catch((error: any) => {
          this.toast.showError(error);
          this.close();
        });
  }

  public selectProject(project: Project.AsObject | GrantedProject.AsObject, type: ProjectType): void {
    if (type === ProjectType.PROJECTTYPE_OWNED) {
      this.project = project as Project.AsObject;
    } else if (type === ProjectType.PROJECTTYPE_GRANTED) {
      this.grantedProject = project as GrantedProject.AsObject;
    }
  }

  public selectGroups(group: Group.AsObject[]): void {
    if (group && group.length) {
      this.groupIds = (group as Group.AsObject[]).map((u) => u.id);
    }
  }

  public selectRoles(roleKeys: string[]): void {
    this.rolesList = roleKeys;
  }

  public next(): void {
    this.currentCreateStep++;
  }

  public previous(): void {
    this.currentCreateStep--;
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
