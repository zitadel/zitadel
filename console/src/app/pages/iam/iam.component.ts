import { Component, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';
import { Features } from 'src/app/proto/generated/zitadel/features_pb';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-iam',
  templateUrl: './iam.component.html',
  styleUrls: ['./iam.component.scss'],
})
export class IamComponent implements OnDestroy {
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  public totalMemberResult: number = 0;
  public membersSubject: BehaviorSubject<Member.AsObject[]> = new BehaviorSubject<Member.AsObject[]>([]);

  public features!: Features.AsObject;

  constructor(
    public adminService: AdminService,
    private dialog: MatDialog,
    private toast: ToastService,
    private breadcrumbService: BreadcrumbService,
    private router: Router,
  ) {
    this.loadMembers();
    this.loadFeatures();
    this.adminService.getDefaultFeatures();

    const breadcrumbs = [
      new Breadcrumb({
        type: BreadcrumbType.IAM,
        name: 'IAM',
        routerLink: ['/iam'],
      }),
    ];
    this.breadcrumbService.setBreadcrumb(breadcrumbs);
  }

  public ngOnDestroy(): void {
    this.breadcrumbService.setBreadcrumb([]);
  }

  public loadMembers(): void {
    this.loadingSubject.next(true);
    from(this.adminService.listIAMMembers(100, 0))
      .pipe(
        map((resp) => {
          if (resp.details?.totalResult) {
            this.totalMemberResult = resp.details.totalResult;
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

  public openAddMember(): void {
    const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
      data: {
        creationType: CreationType.IAM,
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
              return this.adminService.addIAMMember(user.id, roles);
            }),
          )
            .then(() => {
              this.toast.showInfo('IAM.TOAST.MEMBERADDED');
              setTimeout(() => {
                this.loadMembers();
              }, 1000);
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      }
    });
  }

  public showDetail(): void {
    this.router.navigate(['iam/members']);
  }

  public loadFeatures(): void {
    this.loadingSubject.next(true);
    this.adminService.getDefaultFeatures().then((resp) => {
      if (resp.features) {
        this.features = resp.features;
      }
    });
  }
}
