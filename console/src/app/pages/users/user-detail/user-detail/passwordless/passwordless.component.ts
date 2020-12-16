import { Component, Input, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSort } from '@angular/material/sort';
import { MatTable, MatTableDataSource } from '@angular/material/table';
import { BehaviorSubject, Observable } from 'rxjs';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { MFAState, WebAuthNToken } from 'src/app/proto/generated/auth_pb';
import { UserView } from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

export interface WebAuthNOptions {
    challenge: string;
    rp: { name: string, id: string; };
    user: { name: string, id: string, displayName: string; };
    pubKeyCredParams: any;
    authenticatorSelection: { userVerification: string; };
    timeout: number;
    attestation: string;
}

@Component({
    selector: 'app-passwordless',
    templateUrl: './passwordless.component.html',
    styleUrls: ['./passwordless.component.scss'],
})
export class PasswordlessComponent implements OnInit, OnDestroy {
    @Input() private user!: UserView.AsObject;
    public displayedColumns: string[] = ['name', 'state', 'actions'];
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    @ViewChild(MatTable) public table!: MatTable<WebAuthNToken.AsObject>;
    @ViewChild(MatSort) public sort!: MatSort;
    public dataSource!: MatTableDataSource<WebAuthNToken.AsObject>;

    public MFAState: any = MFAState;
    public error: string = '';

    constructor(private service: ManagementService,
        private toast: ToastService,
        private dialog: MatDialog) { }

    public ngOnInit(): void {
        this.getPasswordless();
    }

    public ngOnDestroy(): void {
        this.loadingSubject.complete();
    }

    public getPasswordless(): void {
        this.service.GetPasswordless(this.user.id).then(passwordless => {
            console.log(passwordless.toObject().tokensList);
            this.dataSource = new MatTableDataSource(passwordless.toObject().tokensList);
            this.dataSource.sort = this.sort;
        }).catch(error => {
            this.error = error.message;
        });
    }

    public deletePasswordless(id?: string): void {
        const dialogRef = this.dialog.open(WarnDialogComponent, {
            data: {
                confirmKey: 'ACTIONS.DELETE',
                cancelKey: 'ACTIONS.CANCEL',
                titleKey: 'USER.PASSWORDLESS.DIALOG.DELETE_TITLE',
                descriptionKey: 'USER.PASSWORDLESS.DIALOG.DELETE_DESCRIPTION',
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp && id) {
                this.service.RemovePasswordless(id, this.user.id).then(() => {
                    this.toast.showInfo('USER.TOAST.PASSWORDLESSREMOVED', true);
                    this.getPasswordless();
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        });
    }
}
