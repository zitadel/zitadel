import { Component, OnInit } from '@angular/core';
import { MatButtonToggleChange } from '@angular/material/button-toggle';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';
import { Features } from 'src/app/proto/generated/zitadel/features_pb';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { Org, OrgState } from 'src/app/proto/generated/zitadel/org_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-org-detail',
  templateUrl: './org-detail.component.html',
  styleUrls: ['./org-detail.component.scss'],
})
export class OrgDetailComponent implements OnInit {
  public org!: Org.AsObject;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;

  public OrgState: any = OrgState;
  public ChangeType: any = ChangeType;

  // members
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  public totalMemberResult: number = 0;
  public membersSubject: BehaviorSubject<Member.AsObject[]> = new BehaviorSubject<Member.AsObject[]>([]);
  public features!: Features.AsObject;

  constructor(
    private dialog: MatDialog,
    public mgmtService: ManagementService,
    private toast: ToastService,
    private router: Router,
  ) {}

  public ngOnInit(): void {
    this.getData();
  }

  private async getData(): Promise<void> {
    this.mgmtService
      .getMyOrg()
      .then((resp) => {
        if (resp.org) {
          this.org = resp.org;
          console.log(resp.org);
        }
      })
      .catch((error) => {
        this.toast.showError(error);
      });
    this.loadMembers();
    this.loadFeatures();
  }

  public changeState(event: MatButtonToggleChange | any): void {
    if (event.value === OrgState.ORG_STATE_ACTIVE) {
      this.mgmtService
        .reactivateOrg()
        .then(() => {
          this.toast.showInfo('ORG.TOAST.REACTIVATED', true);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    } else if (event.value === OrgState.ORG_STATE_INACTIVE) {
      this.mgmtService
        .deactivateOrg()
        .then(() => {
          this.toast.showInfo('ORG.TOAST.DEACTIVATED', true);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
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
              this.toast.showError(error);
            });
        }
      }
    });
  }

  public showDetail(): void {
    this.router.navigate(['org/members']);
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

  public loadFeatures(): void {
    this.loadingSubject.next(true);
    this.mgmtService.getFeatures().then((resp) => {
      if (resp.features) {
        this.features = resp.features;
      }
    });
  }
}
