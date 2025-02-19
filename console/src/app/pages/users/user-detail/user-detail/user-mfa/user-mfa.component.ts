import { Component, Input, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { combineLatestWith, defer, EMPTY, Observable, ReplaySubject, Subject, switchMap } from 'rxjs';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { ToastService } from 'src/app/services/toast.service';
import { UserService } from 'src/app/services/user.service';
import { AuthFactor, AuthFactorState, User } from '@zitadel/proto/zitadel/user/v2/user_pb';
import { catchError, filter, map, startWith } from 'rxjs/operators';
import { pairwiseStartWith } from 'src/app/utils/pairwiseStartWith';

export interface MFAItem {
  name: string;
  verified: boolean;
}

type MFAQuery =
  | { state: 'success'; value: MatTableDataSource<AuthFactor>; user: User }
  | { state: 'loading'; value: MatTableDataSource<AuthFactor>; user: User };

@Component({
  selector: 'cnsl-user-mfa',
  templateUrl: './user-mfa.component.html',
  styleUrls: ['./user-mfa.component.scss'],
})
export class UserMfaComponent {
  @Input({ required: true }) public set user(user: User) {
    this.user$.next(user);
  }

  @ViewChild(MatSort) public sort!: MatSort;
  public dataSource = new MatTableDataSource<AuthFactor>([]);

  public displayedColumns: string[] = ['type', 'name', 'state', 'actions'];
  private user$ = new ReplaySubject<User>(1);
  public mfaQuery$: Observable<MFAQuery>;
  public refresh$ = new Subject<true>();
  public AuthFactorState = AuthFactorState;

  constructor(
    private readonly dialog: MatDialog,
    private readonly toast: ToastService,
    private readonly userService: UserService,
  ) {
    this.mfaQuery$ = this.user$.pipe(
      combineLatestWith(this.refresh$.pipe(startWith(true))),
      switchMap(([user]) => this.listAuthenticationFactors(user)),
      pairwiseStartWith(undefined),
      map(([prev, curr]) => {
        if (prev?.state === 'success' && curr.state === 'loading') {
          return { ...prev, state: 'loading' } as const;
        }
        return curr;
      }),
      catchError((error) => {
        this.toast.showError(error);
        return EMPTY;
      }),
    );
  }

  private listAuthenticationFactors(user: User): Observable<MFAQuery> {
    return defer(() => this.userService.listAuthenticationFactors({ userId: user.userId })).pipe(
      map(
        ({ result }) =>
          ({
            state: 'success',
            value: new MatTableDataSource<AuthFactor>(result),
            user,
          }) as const,
      ),
      startWith({
        state: 'loading',
        value: new MatTableDataSource<AuthFactor>([]),
        user,
      } as const),
    );
  }

  private async removeTOTP(user: User) {
    await this.userService.removeTOTP(user.userId);
    return ['USER.TOAST.OTPREMOVED', 'otp'] as const;
  }

  private async removeU2F(user: User, u2fId: string) {
    await this.userService.removeU2F(user.userId, u2fId);
    return ['USER.TOAST.U2FREMOVED', 'u2f'] as const;
  }

  private async removeOTPEmail(user: User) {
    await this.userService.removeOTPEmail(user.userId);
    return ['USER.TOAST.OTPREMOVED', 'otpEmail'] as const;
  }

  private async removeOTPSMS(user: User) {
    await this.userService.removeOTPSMS(user.userId);
    return ['USER.TOAST.OTPREMOVED', 'otpSms'] as const;
  }

  public deleteMFA(user: User, factor: AuthFactor): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'USER.MFA.DIALOG.MFA_DELETE_TITLE',
        descriptionKey: 'USER.MFA.DIALOG.MFA_DELETE_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef
      .afterClosed()
      .pipe(
        filter(Boolean),
        switchMap(() => {
          switch (factor.type.case) {
            case 'otp':
              return this.removeTOTP(user);
            case 'u2f':
              return this.removeU2F(user, factor.type.value.id);
            case 'otpEmail':
              return this.removeOTPEmail(user);
            case 'otpSms':
              return this.removeOTPSMS(user);
            default:
              throw new Error('Unknown MFA type');
          }
        }),
      )
      .subscribe({
        next: ([translation, caseId]) => {
          this.toast.showInfo(translation, true);
          const index = this.dataSource.data.findIndex((mfa) => mfa.type.case === caseId);
          if (index > -1) {
            this.dataSource.data.splice(index, 1);
          }
          this.refresh$.next(true);
        },
        error: (error) => this.toast.showError(error),
      });
  }
}
