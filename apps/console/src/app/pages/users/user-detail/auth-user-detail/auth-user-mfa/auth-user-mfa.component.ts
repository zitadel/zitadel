import { Component, Input, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSort } from '@angular/material/sort';
import { MatTable, MatTableDataSource } from '@angular/material/table';
import { BehaviorSubject, Observable } from 'rxjs';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { AuthFactorState } from 'src/app/proto/generated/zitadel/user_pb';
import { NewAuthService } from 'src/app/services/new-auth.service';
import { ToastService } from 'src/app/services/toast.service';
import { AddAuthFactorDialogData, AuthFactorDialogComponent } from '../auth-factor-dialog/auth-factor-dialog.component';
import { AuthFactor } from '@zitadel/proto/zitadel/user_pb';
import { SecondFactorType } from '@zitadel/proto/zitadel/policy_pb';
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
  selector: 'cnsl-auth-user-mfa',
  templateUrl: './auth-user-mfa.component.html',
  styleUrls: ['./auth-user-mfa.component.scss'],
})
export class AuthUserMfaComponent implements OnInit, OnDestroy {
  public displayedColumns: string[] = ['name', 'type', 'state', 'actions'];
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();

  @ViewChild(MatTable) public table!: MatTable<AuthFactor>;
  @ViewChild(MatSort) public sort!: MatSort;
  @Input() public phoneVerified: boolean = false;
  public AuthFactorState: any = AuthFactorState;
  public dataSource: MatTableDataSource<AuthFactor> = new MatTableDataSource<AuthFactor>([]);

  protected error: string = '';

  protected otpAvailable$ = new BehaviorSubject<boolean>(false);
  protected u2fAvailable$ = new BehaviorSubject<boolean>(false);
  protected otpSmsAvailable$ = new BehaviorSubject<boolean>(false);
  protected otpEmailAvailable$ = new BehaviorSubject<boolean>(false);
  protected otpDisabled$ = new BehaviorSubject<boolean>(true);
  protected otpSmsDisabled$ = new BehaviorSubject<boolean>(true);
  protected otpEmailDisabled$ = new BehaviorSubject<boolean>(true);

  constructor(
    private readonly service: NewAuthService,
    private readonly toast: ToastService,
    private readonly dialog: MatDialog,
  ) {}

  public ngOnInit(): void {
    this.getMFAs();
    this.applyOrgPolicy();
  }

  public ngOnDestroy(): void {
    this.loadingSubject.complete();
  }

  public addAuthFactor(): void {
    const data: AddAuthFactorDialogData = {
      otp$: this.otpAvailable$,
      u2f$: this.u2fAvailable$,
      otpSms$: this.otpSmsAvailable$,
      otpEmail$: this.otpEmailAvailable$,
      otpDisabled$: this.otpDisabled$,
      otpSmsDisabled$: this.otpSmsDisabled$,
      otpEmailDisabled$: this.otpEmailDisabled$,
      phoneVerified: this.phoneVerified,
    } as const;

    const dialogRef = this.dialog.open(AuthFactorDialogComponent, {
      data: data,
    });

    dialogRef.afterClosed().subscribe(() => {
      this.getMFAs();
    });
  }

  public getMFAs(): void {
    this.service
      .listMyMultiFactors()
      .then((mfas) => {
        const list: AuthFactor[] = mfas.result;
        this.dataSource = new MatTableDataSource(list);
        this.dataSource.sort = this.sort;

        this.disableAuthFactor(list, 'otp', this.otpDisabled$);
        this.disableAuthFactor(list, 'otpSms', this.otpSmsDisabled$);
        this.disableAuthFactor(list, 'otpEmail', this.otpEmailDisabled$);
      })
      .catch((error) => {
        this.error = error.message;
      });
  }

  public applyOrgPolicy(): void {
    this.service.getMyLoginPolicy().then((resp) => {
      if (resp && resp.policy) {
        const secondFactors = resp.policy?.secondFactors;
        this.displayAuthFactorBasedOnPolicy(secondFactors, SecondFactorType.OTP, this.otpAvailable$);
        this.displayAuthFactorBasedOnPolicy(secondFactors, SecondFactorType.U2F, this.u2fAvailable$);
        this.displayAuthFactorBasedOnPolicy(secondFactors, SecondFactorType.OTP_EMAIL, this.otpEmailAvailable$);
        this.displayAuthFactorBasedOnPolicy(secondFactors, SecondFactorType.OTP_SMS, this.otpSmsAvailable$);
      }
    });
  }

  public deleteMFA(factor: AuthFactor): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'USER.MFA.DIALOG.MFA_DELETE_TITLE',
        descriptionKey: 'USER.MFA.DIALOG.MFA_DELETE_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        if (factor.type.case === 'otp') {
          this.service
            .removeMyMultiFactorOTP()
            .then(() => {
              this.toast.showInfo('USER.TOAST.OTPREMOVED', true);

              this.cleanupList();
              this.getMFAs();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        } else if (factor.type.case === 'u2f') {
          this.service
            .removeMyMultiFactorU2F(factor.type.value.id)
            .then(() => {
              this.toast.showInfo('USER.TOAST.U2FREMOVED', true);

              this.cleanupList();
              this.getMFAs();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        } else if (factor.type.case === 'otpEmail') {
          this.service
            .removeMyAuthFactorOTPEmail()
            .then(() => {
              this.toast.showInfo('USER.TOAST.OTPREMOVED', true);

              this.cleanupList();
              this.getMFAs();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        } else if (factor.type.case === 'otpSms') {
          this.service
            .removeMyAuthFactorOTPSMS()
            .then(() => {
              this.toast.showInfo('USER.TOAST.OTPREMOVED', true);

              this.cleanupList();
              this.getMFAs();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      }
    });
  }

  private cleanupList(): void {
    this.dataSource.data = this.dataSource.data.filter((mfa: AuthFactor) => {
      return mfa.type.case;
    });
  }

  private disableAuthFactor(mfas: AuthFactor[], key: string, subject: BehaviorSubject<boolean>): void {
    subject.next(mfas.some((mfa) => mfa.type.case === key));
  }

  private displayAuthFactorBasedOnPolicy(
    factors: SecondFactorType[],
    factor: SecondFactorType,
    subject: BehaviorSubject<boolean>,
  ): void {
    subject.next(factors.some((f) => f === factor));
  }
}
