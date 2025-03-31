import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTableDataSource } from '@angular/material/table';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable, of } from 'rxjs';
import { enterAnimations } from 'src/app/animations';
import { Timestamp } from 'src/app/proto/generated/google/protobuf/timestamp_pb';
import { Group, GroupState } from 'src/app/proto/generated/zitadel/group_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { WarnDialogComponent } from '../warn-dialog/warn-dialog.component';
import { GroupCreateDialogComponent } from '../add-group-dialog/group-create-dialog.component';

@Component({
    selector: 'cnsl-user-groups',
    templateUrl: './user-groups.component.html',
    styleUrls: ['./user-groups.component.scss'],
    animations: [enterAnimations],
})
export class UserGroupsComponent implements OnInit {
    @Input() refreshOnPreviousRoutes: string[] = [];
    @Input() public canWrite$: Observable<boolean> = of(false);
    @Input() public canDelete$: Observable<boolean> = of(false);
    @Input() userId: string = '';

    public selection: SelectionModel<Group.AsObject> = new SelectionModel<Group.AsObject>(true, []);
    public dataSource: MatTableDataSource<Group.AsObject> = new MatTableDataSource<Group.AsObject>();
    public totalResult: number = 0;
    public viewTimestamp!: Timestamp.AsObject;
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable < boolean > = this.loadingSubject.asObservable();
    public GroupState: any = GroupState;
    public changePage: EventEmitter<void> = new EventEmitter();
    @Input() public displayedColumnsHuman: string[] = [
        'select',
        'name',
        'state',
        'creationDate',
        'changeDate',
        'actions',
    ];

    constructor(
        public translate: TranslateService,
        private groupService: ManagementService,
        private toast: ToastService,
        private dialog: MatDialog,
      ) {
        
      }

    public masterToggle(): void {
        this.isAllSelected() ? this.selection.clear() : this.dataSource.data.forEach((row) => this.selection.select(row));
    }

    ngOnInit(): void {
        this.changePage.emit();
        this.getData();
    }

    private getData(): void {
        this.loadingSubject.next(true);
        this.groupService.getGroupByUserID(this.userId).then((resp) => {
            if (resp.details?.totalResult) {
                this.totalResult = resp.details?.totalResult;
            } else {
                this.totalResult = 0;
            }
            if (resp.details?.viewTimestamp) {
                this.viewTimestamp = resp.details?.viewTimestamp;
            }
            this.dataSource.data = resp.resultList;
            this.loadingSubject.next(false);
        }).catch((error) => {
            this.toast.showError(error);
            this.loadingSubject.next(false);
        });
    }

    public isAllSelected(): boolean {
        const numSelected = this.selection.selected.length;
        const numRows = this.dataSource.data.length;
        return numSelected === numRows;
    }

    public deleteGroup(group: Group.AsObject): void {
        const dialogRef = this.dialog.open(WarnDialogComponent, {
            data: {
                confirmKey: 'ACTIONS.DELETE',
                cancelKey: 'ACTIONS.CANCEL',
                titleKey: 'USER.GROUPS.DELETE_TITLE',
                descriptionKey: 'USER.GROUPS.DELETE_DESCRIPTION',
            },
            width: '400px',
        });
        dialogRef.afterClosed().subscribe((resp) => {
            if (resp) {
                this.groupService
                .removeGroupMember(group.id, this.userId)
                .then(() => {
                    this.toast.showInfo('GROUP.TOAST.MEMBERREMOVED', true);
                    this.getData();
                })
                .catch((error) => {
                    this.toast.showError(error);
                });
            }
        });
    }

    public refreshPage(): void {
        this.getData();
    }

    public openAddGroup(): void {
        const dialogRef = this.dialog.open(GroupCreateDialogComponent, {width: '400px'});
        dialogRef.afterClosed().subscribe((resp) => {
            if (resp) {
                const groups: string[] = resp.groups;
                if (groups && groups.length) {
                    Promise.all(
                        groups.map((group) => {
                            return this.groupService.addGroupMember(group, this.userId);
                        }),
                    )
                    .then(() => {
                        setTimeout(() => {
                            this.changePage.emit();
                            this.getData();
                        }, 1000);
                        this.toast.showInfo('GROUP.TOAST.MEMBERSADDED', true);
                    })
                    .catch((error) => {
                        this.changePage.emit();
                        this.toast.showError(error);
                    });
                }
            }
        });
    }
}