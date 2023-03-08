import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { MatLegacyTable as MatTable, MatLegacyTableDataSource as MatTableDataSource } from '@angular/material/legacy-table';
import { MatSort } from '@angular/material/sort';
import { BehaviorSubject, Observable } from 'rxjs';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { AuthFactor, AuthFactorState } from 'src/app/proto/generated/zitadel/user_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

import { AuthFactorDialogComponent } from '../auth-factor-dialog/auth-factor-dialog.component';

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

  @ViewChild(MatTable) public table!: MatTable<AuthFactor.AsObject>;
  @ViewChild(MatSort) public sort!: MatSort;
  public dataSource: MatTableDataSource<AuthFactor.AsObject> = new MatTableDataSource<AuthFactor.AsObject>([]);

  public AuthFactorState: any = AuthFactorState;

  public error: string = '';
  public otpDisabled$ = new BehaviorSubject<boolean>(true);

  constructor(private service: GrpcAuthService, private toast: ToastService, private dialog: MatDialog) {}

  public ngOnInit(): void {
    this.getMFAs();
  }

  public ngOnDestroy(): void {
    this.loadingSubject.complete();
  }

  public addAuthFactor(): void {
    const dialogRef = this.dialog.open(AuthFactorDialogComponent, {
      data: {
        otpDisabled$: this.otpDisabled$,
      },
    });

    dialogRef.afterClosed().subscribe(() => {
      this.getMFAs();
    });
  }

  public getMFAs(): void {
    this.service
      .listMyMultiFactors()
      .then((mfas) => {
        const list = mfas.resultList;
        this.dataSource = new MatTableDataSource(list);
        this.dataSource.sort = this.sort;

        const index = list.findIndex((mfa) => mfa.otp);
        if (index === -1) {
          this.otpDisabled$.next(false);
        }
      })
      .catch((error) => {
        this.error = error.message;
      });
  }

  public deleteMFA(factor: AuthFactor.AsObject): void {
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
        if (factor.otp) {
          this.service
            .removeMyMultiFactorOTP()
            .then(() => {
              this.toast.showInfo('USER.TOAST.OTPREMOVED', true);

              const index = this.dataSource.data.findIndex((mfa) => !!mfa.otp);
              if (index > -1) {
                this.dataSource.data.splice(index, 1);
              }
              this.getMFAs();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        } else if (factor.u2f) {
          this.service
            .removeMyMultiFactorU2F(factor.u2f.id)
            .then(() => {
              this.toast.showInfo('USER.TOAST.U2FREMOVED', true);

              const index = this.dataSource.data.findIndex((mfa) => !!mfa.u2f);
              if (index > -1) {
                this.dataSource.data.splice(index, 1);
              }
              this.getMFAs();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      }
    });
  }
}
