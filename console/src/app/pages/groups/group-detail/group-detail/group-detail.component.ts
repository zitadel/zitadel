import { MediaMatcher } from '@angular/cdk/layout';
import { Location } from '@angular/common';
import { PageEvent } from '@angular/material/paginator';
import { Component, EventEmitter, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { take } from 'rxjs/operators';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { SidenavSetting } from 'src/app/modules/sidenav/sidenav.component';
import { GroupGrantContext } from 'src/app/modules/group-grants/group-grants-datasource';
import { Group, GroupState } from 'src/app/proto/generated/zitadel/group_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { LanguagesService } from '../../../../services/languages.service';
import { NameDialogComponent } from 'src/app/modules/name-dialog/name-dialog.component';
import { GroupMemberCreateDialogComponent } from 'src/app/modules/add-group-member-dialog/group-member-create-dialog.component';
import { GroupMembersDataSource } from './group-members-datasource';

@Component({
  selector: 'cnsl-group-detail',
  templateUrl: './group-detail.component.html',
  styleUrls: ['./group-detail.component.scss'],
})
export class GroupDetailComponent implements OnInit {
  public group!: Group.AsObject;
  public groupId: string = '';

  public loading: boolean = true;

  public GroupState: any = GroupState;
  public ChangeType: any = ChangeType;

  public changePage: EventEmitter<void> = new EventEmitter();
  public settingsList: SidenavSetting[] = [
    { id: 'members', i18nKey: 'GROUP.SETTINGS.MEMBERS' },
    { id: 'grants', i18nKey: 'GROUP.SETTINGS.GROUPGRANTS' },
  ];
  public currentSetting: string | undefined = this.settingsList[0].id;
  public GROUPGRANTCONTEXT: GroupGrantContext = GroupGrantContext.GROUP;

  public error: string = '';

  public changePageFactory!: Function;
  public dataSource!: GroupMembersDataSource;
  public groupName: string = '';
  public INITIALPAGESIZE: number = 25;

  constructor(
    public translate: TranslateService,
    private route: ActivatedRoute,
    private toast: ToastService,
    public mgmtGroupService: ManagementService,
    private _location: Location,
    private dialog: MatDialog,
    private router: Router,
    activatedRoute: ActivatedRoute,
    private mediaMatcher: MediaMatcher,
    public langSvc: LanguagesService,
    breadcrumbService: BreadcrumbService,
  ) {
    activatedRoute.queryParams.pipe(take(1)).subscribe((params: Params) => {
      const { key } = params;
      if (key) {
        this.currentSetting = key;
      }
    });
    breadcrumbService.setBreadcrumb([
      new Breadcrumb({
        type: BreadcrumbType.ORG,
        routerLink: ['/org'],
      }),
    ]);

    this.route.params.subscribe((params) => {
      this.groupId = params['id'];
      this.loadMembers();
    });
  }


  refreshGroup(): void {
    this.changePage.emit();
    this.route.params.pipe(take(1)).subscribe((params) => {
      this.loading = true;
      const { id } = params;
      this.mgmtGroupService
        .getGroupByID(id)
        .then((resp) => {
          this.loading = false;
          if (resp.group) {
            this.group = resp.group;
          }
        })
        .catch((err) => {
          this.error = err.message ?? '';
          this.loading = false;
          this.toast.showError(err);
        });
    });
  }

  public ngOnInit(): void {
    const groupId = this.route.snapshot.paramMap.get('id');
    this.refreshGroup();
  }

  public changeState(newState: GroupState): void {
    if (newState === GroupState.GROUP_STATE_ACTIVE) {
      this.mgmtGroupService
        .reactivateGroup(this.group.id)
        .then(() => {
          this.toast.showInfo('GROUP.TOAST.REACTIVATED', true);
          this.group.state = newState;
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    } else if (newState === GroupState.GROUP_STATE_INACTIVE) {
      this.mgmtGroupService
        .deactivateGroup(this.group.id)
        .then(() => {
          this.toast.showInfo('GROUP.TOAST.DEACTIVATED', true);
          this.group.state = newState;
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public navigateBack(): void {
    this._location.back();
  }

  public deleteGroup(): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'GROUP.DIALOG.DELETE_TITLE',
        descriptionKey: 'GROUP.DIALOG.DELETE_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.mgmtGroupService
          .removeGroup(this.group.id)
          .then(() => {
            const params: Params = {
              deferredReload: true,
            };
            this.router.navigate(['/groups'], { queryParams: params });
            this.toast.showInfo('GROUP.TOAST.DELETED', true);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public openNameDialog(): void {
    const dialogRef = this.dialog.open(NameDialogComponent, {
      data: {
        name: this.group?.name,
        titleKey: 'APP.NAMEDIALOG.TITLE',
        descKey: 'APP.NAMEDIALOG.DESCRIPTION',
        labelKey: 'APP.NAMEDIALOG.NAME',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((name) => {
      if (name) {
        this.group!.name = name;
        this.saveGroup();
      }
    });
  }

  public saveGroup(): void {
    if (this.group) {
      this.mgmtGroupService
        .updateGroup(this.group.id, this.group.name)
        .then(() => {
          this.toast.showInfo('APP.TOAST.UPDATED', true);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public openAddMember(): void {
    const dialogRef = this.dialog.open(GroupMemberCreateDialogComponent, {width: '400px'});
    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        const users: User.AsObject[] = resp.users;
        if (users && users.length) {
          Promise.all(
            users.map((user) => {
              return this.mgmtGroupService.addGroupMember((this.group as Group.AsObject).id, user.id);
            }),
          )
            .then(() => {
              setTimeout(() => {
                this.changePage.emit();
              }, 1000);
              this.toast.showInfo('PROJECT.TOAST.MEMBERSADDED', true);
            })
            .catch((error) => {
              this.changePage.emit();
              this.toast.showError(error);
            });
        }
      }
    });
  }

  public removeGroupMember(member: Member.AsObject | Member.AsObject): void {
    this.mgmtGroupService
      .removeGroupMember((this.group as Group.AsObject).id, member.userId)
      .then(() => {
        setTimeout(() => {
          this.changePage.emit();
        }, 1000);
        this.toast.showInfo('PROJECT.TOAST.MEMBERREMOVED', true);
      })
      .catch((error) => {
        this.toast.showError(error);
        this.changePage.emit();
      });
  }

  public loadMembers(): Promise<any> {
    return this.mgmtGroupService.getGroupByID(this.groupId).then((resp) => {
      if (resp.group) {
        this.group = resp.group;
        this.groupName = this.group.name;
        this.dataSource = new GroupMembersDataSource(this.mgmtGroupService);
        this.dataSource.loadMembers(this.group.id, 0, this.INITIALPAGESIZE);

        this.changePageFactory = (event?: PageEvent) => {
          return this.dataSource.loadMembers(
            (this.group as Group.AsObject).id,
            event?.pageIndex ?? 0,
            event?.pageSize ?? this.INITIALPAGESIZE,
          );
        };
      }
    });
  }
}
