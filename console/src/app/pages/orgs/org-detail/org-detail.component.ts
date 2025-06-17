import { ChangeDetectorRef, Component, effect, OnInit, signal } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { BehaviorSubject, from, lastValueFrom, Observable, of } from 'rxjs';
import { catchError, distinctUntilChanged, finalize, map } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { MetadataDialogComponent } from 'src/app/modules/metadata/metadata-dialog/metadata-dialog.component';
import { NameDialogComponent } from 'src/app/modules/name-dialog/name-dialog.component';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { Metadata } from 'src/app/proto/generated/zitadel/metadata_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { NewOrganizationService } from '../../../services/new-organization.service';
import { injectMutation } from '@tanstack/angular-query-experimental';
import { Organization, OrganizationState } from '@zitadel/proto/zitadel/org/v2/org_pb';
import { toObservable } from '@angular/core/rxjs-interop';

@Component({
  selector: 'cnsl-org-detail',
  templateUrl: './org-detail.component.html',
  styleUrls: ['./org-detail.component.scss'],
})
export class OrgDetailComponent implements OnInit {
  public PolicyComponentServiceType: any = PolicyComponentServiceType;

  public OrganizationState = OrganizationState;
  public ChangeType: any = ChangeType;

  public metadata: Metadata.AsObject[] = [];
  public loadingMetadata: boolean = true;

  // members
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  public totalMemberResult: number = 0;
  public membersSubject: BehaviorSubject<Member.AsObject[]> = new BehaviorSubject<Member.AsObject[]>([]);

  public InfoSectionType: any = InfoSectionType;

  protected readonly orgQuery = this.newOrganizationService.activeOrganizationQuery();
  private readonly reactivateOrgMutation = injectMutation(this.newOrganizationService.reactivateOrgMutationOptions);
  private readonly deactivateOrgMutation = injectMutation(this.newOrganizationService.deactivateOrgMutationOptions);
  private readonly deleteOrgMutation = injectMutation(this.newOrganizationService.deleteOrgMutationOptions);
  private readonly renameOrgMutation = injectMutation(this.newOrganizationService.renameOrgMutationOptions);

  protected reloadChanges = signal(true);

  constructor(
    private readonly dialog: MatDialog,
    private readonly mgmtService: ManagementService,
    private readonly toast: ToastService,
    private readonly router: Router,
    private readonly newOrganizationService: NewOrganizationService,
    breadcrumbService: BreadcrumbService,
    cdr: ChangeDetectorRef,
  ) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);

    effect(() => {
      const orgId = this.newOrganizationService.orgId();
      if (!orgId) {
        return;
      }
      this.loadMembers();
      this.loadMetadata();
    });

    // force rerender changes because it is not reactive to orgId changes
    toObservable(this.newOrganizationService.orgId).subscribe(() => {
      this.reloadChanges.set(false);
      cdr.detectChanges();
      this.reloadChanges.set(true);
    });
  }

  public ngOnInit(): void {
    this.loadMembers();
    this.loadMetadata();
  }

  public async changeState(newState: OrganizationState) {
    if (newState === OrganizationState.ACTIVE) {
      const dialogRef = this.dialog.open(WarnDialogComponent, {
        data: {
          confirmKey: 'ACTIONS.REACTIVATE',
          cancelKey: 'ACTIONS.CANCEL',
          titleKey: 'ORG.DIALOG.REACTIVATE.TITLE',
          descriptionKey: 'ORG.DIALOG.REACTIVATE.DESCRIPTION',
        },
        width: '400px',
      });
      const resp = await lastValueFrom(dialogRef.afterClosed());
      if (!resp) {
        return;
      }
      try {
        await this.reactivateOrgMutation.mutateAsync();
        this.toast.showInfo('ORG.TOAST.REACTIVATED', true);
      } catch (error) {
        this.toast.showError(error);
      }
      return;
    }

    if (newState === OrganizationState.INACTIVE) {
      const dialogRef = this.dialog.open(WarnDialogComponent, {
        data: {
          confirmKey: 'ACTIONS.DEACTIVATE',
          cancelKey: 'ACTIONS.CANCEL',
          titleKey: 'ORG.DIALOG.DEACTIVATE.TITLE',
          descriptionKey: 'ORG.DIALOG.DEACTIVATE.DESCRIPTION',
        },
        width: '400px',
      });

      const resp = await lastValueFrom(dialogRef.afterClosed());
      if (!resp) {
        return;
      }
      try {
        await this.deactivateOrgMutation.mutateAsync();
        this.toast.showInfo('ORG.TOAST.DEACTIVATED', true);
      } catch (error) {
        this.toast.showError(error);
      }
    }
  }

  public async deleteOrg(org: Organization) {
    const mgmtUserData = {
      confirmKey: 'ACTIONS.DELETE',
      cancelKey: 'ACTIONS.CANCEL',
      titleKey: 'ORG.DIALOG.DELETE.TITLE',
      warnSectionKey: 'ORG.DIALOG.DELETE.DESCRIPTION',
      hintKey: 'ORG.DIALOG.DELETE.TYPENAME',
      hintParam: 'ORG.DIALOG.DELETE.DESCRIPTION',
      confirmationKey: 'ORG.DIALOG.DELETE.ORGNAME',
      confirmation: org.name,
    };

    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: mgmtUserData,
      width: '400px',
    });

    if (!(await lastValueFrom(dialogRef.afterClosed()))) {
      return;
    }

    try {
      await this.deleteOrgMutation.mutateAsync();
      await this.router.navigate(['/orgs']);
    } catch (error) {
      this.toast.showError(error);
    }
  }

  public openAddMember(): void {
    const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
      data: {
        creationType: CreationType.ORG,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        const users: User.AsObject[] = resp.users;
        const roles: string[] = resp.roles;

        if (users && users.length && roles && roles.length) {
          Promise.all(
            users.map((user) => {
              return this.mgmtService.addOrgMember(user.id, roles);
            }),
          )
            .then(() => {
              this.toast.showInfo('ORG.TOAST.MEMBERADDED', true);
              setTimeout(() => {
                this.loadMembers();
              }, 1000);
            })
            .catch((error) => {
              setTimeout(() => {
                this.loadMembers();
              }, 1000);
              this.toast.showError(error);
            });
        }
      }
    });
  }

  public showDetail() {
    return this.router.navigate(['org/members']);
  }

  public loadMembers(): void {
    this.loadingSubject.next(true);
    from(this.mgmtService.listOrgMembers(100, 0))
      .pipe(
        map((resp) => {
          if (resp.details?.totalResult) {
            this.totalMemberResult = resp.details?.totalResult;
          } else {
            this.totalMemberResult = 0;
          }

          return resp.resultList;
        }),
        catchError(() => of([])),
        finalize(() => this.loadingSubject.next(false)),
      )
      .subscribe((members) => {
        this.membersSubject.next(members);
      });
  }

  public loadMetadata(): Promise<any> | void {
    this.loadingMetadata = true;
    return this.mgmtService
      .listOrgMetadata()
      .then((resp) => {
        this.loadingMetadata = false;
        const decoder = new TextDecoder();
        this.metadata = resp.resultList.map(({ key, value }) => {
          return {
            key,
            value: atob(typeof value === 'string' ? value : decoder.decode(value)),
          };
        });
      })
      .catch((error) => {
        this.loadingMetadata = false;
        this.toast.showError(error);
      });
  }

  public editMetadata(): void {
    const setFcn = (key: string, value: string): Promise<any> => this.mgmtService.setOrgMetadata(key, btoa(value));
    const removeFcn = (key: string): Promise<any> => this.mgmtService.removeOrgMetadata(key);

    const dialogRef = this.dialog.open(MetadataDialogComponent, {
      data: {
        metadata: this.metadata,
        setFcn: setFcn,
        removeFcn: removeFcn,
      },
    });

    dialogRef.afterClosed().subscribe(() => {
      this.loadMetadata();
    });
  }

  public async renameOrg(org: Organization): Promise<void> {
    const dialogRef = this.dialog.open(NameDialogComponent, {
      data: {
        name: org.name,
        titleKey: 'ORG.PAGES.RENAME.TITLE',
        descKey: 'ORG.PAGES.RENAME.DESCRIPTION',
        labelKey: 'ORG.PAGES.NAME',
      },
      width: '400px',
    });

    const name = await lastValueFrom(dialogRef.afterClosed());
    if (org.name === name) {
      return;
    }

    try {
      await this.renameOrgMutation.mutateAsync(name);
      this.toast.showInfo('ORG.TOAST.UPDATED', true);
      const resp = await this.mgmtService.getMyOrg();
      if (resp.org) {
        await this.newOrganizationService.setOrgId(resp.org.id);
      }
    } catch (error) {
      this.toast.showError(error);
    }
  }
}
