import { Component, Input, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { BehaviorSubject, Observable, switchMap } from 'rxjs';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { ToastService } from 'src/app/services/toast.service';
import { AuthFactorState, Passkey, User } from '@zitadel/proto/zitadel/user/v2/user_pb';
import { UserService } from 'src/app/services/user.service';
import { filter } from 'rxjs/operators';

export interface WebAuthNOptions {
  challenge: string;
  rp: { name: string; id: string };
  user: { name: string; id: string; displayName: string };
  pubKeyCredParams: any;
  authenticatorSelection: { userVerification: string };
  timeout: number;
  attestation: string;
}

@Component({
  selector: 'cnsl-passwordless',
  templateUrl: './passwordless.component.html',
  styleUrls: ['./passwordless.component.scss'],
})
export class PasswordlessComponent implements OnInit, OnDestroy {
  @Input({ required: true }) public user!: User;
  @Input() public disabled: boolean = true;
  public displayedColumns: string[] = ['name', 'state', 'actions'];
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();

  @ViewChild(MatSort) public sort!: MatSort;
  public dataSource: MatTableDataSource<Passkey> = new MatTableDataSource<Passkey>([]);

  public AuthFactorState = AuthFactorState;
  public error: string = '';

  constructor(
    private toast: ToastService,
    private dialog: MatDialog,
    private userService: UserService,
  ) {}

  public ngOnInit(): void {
    this.getPasswordless();
  }

  ngOnDestroy(): void {
    this.loadingSubject.complete();
  }

  public getPasswordless(): void {
    this.userService
      .listPasskeys({ userId: this.user.userId })
      .then((passwordless) => {
        this.dataSource = new MatTableDataSource(passwordless.result);
        this.dataSource.sort = this.sort;
      })
      .catch((error) => {
        this.error = error.message;
      });
  }

  public deletePasswordless(passkeyId: string): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'USER.PASSWORDLESS.DIALOG.DELETE_TITLE',
        descriptionKey: 'USER.PASSWORDLESS.DIALOG.DELETE_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef
      .afterClosed()
      .pipe(
        filter(Boolean),
        switchMap(() => this.userService.removePasskeys({ userId: this.user.userId, passkeyId })),
      )
      .subscribe({
        next: () => {
          this.toast.showInfo('USER.TOAST.PASSWORDLESSREMOVED', true);
          this.getPasswordless();
        },
        error: (error) => {
          this.toast.showError(error);
        },
      });
  }

  public sendPasswordlessRegistration(): void {
    this.userService
      .createPasskeyRegistrationLink({ userId: this.user.userId, medium: { case: 'sendLink', value: {} } })
      .then(() => {
        this.toast.showInfo('USER.TOAST.PASSWORDLESSREGISTRATIONSENT', true);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }
}
