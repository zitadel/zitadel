import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSort } from '@angular/material/sort';
import { MatTable, MatTableDataSource } from '@angular/material/table';
import { BehaviorSubject, Observable } from 'rxjs';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { AuthFactorState, WebAuthNToken } from 'src/app/proto/generated/zitadel/user_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

import { EnvironmentService } from 'src/app/services/environment.service';
import { UserService } from 'src/app/services/user.service';
import { _base64ToArrayBuffer } from '../../u2f-util';
import { DialogPasswordlessComponent } from './dialog-passwordless/dialog-passwordless.component';

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
  selector: 'cnsl-auth-passwordless',
  templateUrl: './auth-passwordless.component.html',
  styleUrls: ['./auth-passwordless.component.scss'],
  standalone: false,
})
export class AuthPasswordlessComponent implements OnInit, OnDestroy {
  public displayedColumns: string[] = ['name', 'state', 'actions'];
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();

  @ViewChild(MatTable) public table!: MatTable<WebAuthNToken.AsObject>;
  @ViewChild(MatSort) public sort!: MatSort;
  public dataSource: MatTableDataSource<WebAuthNToken.AsObject> = new MatTableDataSource<WebAuthNToken.AsObject>([]);

  public AuthFactorState: any = AuthFactorState;
  public error: string = '';

  // WebAuthn Relying Party ID and login app base URL, resolved from the runtime environment
  // (falling back to the current host/origin). The RP ID must match the login app's so that
  // passkeys registered here can be used at login.
  private rpId: string = window.location.hostname;

  constructor(
    private service: GrpcAuthService,
    private userService: UserService,
    private envService: EnvironmentService,
    private toast: ToastService,
    private dialog: MatDialog,
  ) {
    this.envService.env.subscribe((env) => {
      if (env.webauthn_rp_id) {
        this.rpId = env.webauthn_rp_id;
      }
    });
  }

  public ngOnInit(): void {
    this.getPasswordless();
  }

  public ngOnDestroy(): void {
    this.loadingSubject.complete();
  }

  public async addPasswordless(): Promise<void> {
    const userId = this.userService.userId();
    if (!userId) {
      this.toast.showError('USER.PASSWORDLESS.U2F_ERROR', false, true);
      return;
    }

    try {
      // 1. obtain a registration code, 2. start the passkey registration with our RP ID
      const link = await this.userService.createPasskeyRegistrationLink({
        userId,
        medium: { case: 'returnCode', value: {} },
      });
      const resp = await this.userService.registerPasskey({
        userId,
        code: link.code,
        domain: this.rpId,
      });

      const credOptions = resp.publicKeyCredentialCreationOptions as unknown as CredentialCreationOptions;
      if (!credOptions?.publicKey?.challenge) {
        this.toast.showError('USER.PASSWORDLESS.U2F_ERROR', false, true);
        return;
      }

      credOptions.publicKey.challenge = _base64ToArrayBuffer(credOptions.publicKey.challenge as any);
      credOptions.publicKey.user.id = _base64ToArrayBuffer(credOptions.publicKey.user.id as any);
      if (credOptions.publicKey.excludeCredentials) {
        credOptions.publicKey.excludeCredentials.map((cred) => {
          cred.id = _base64ToArrayBuffer(cred.id as any);
          return cred;
        });
      }

      const dialogRef = this.dialog.open(DialogPasswordlessComponent, {
        width: '400px',
        data: {
          credOptions,
          passkeyId: resp.passkeyId,
        },
      });

      dialogRef.afterClosed().subscribe(() => {
        setTimeout(() => {
          this.getPasswordless();
        }, 1000);
      });
    } catch (error) {
      this.toast.showError(error);
    }
  }

  public getPasswordless(): void {
    this.service
      .listMyPasswordless()
      .then((passwordless) => {
        this.dataSource = new MatTableDataSource(passwordless.resultList);
        this.dataSource.sort = this.sort;
      })
      .catch((error) => {
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

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp && id) {
        this.service
          .removeMyPasswordless(id)
          .then(() => {
            this.toast.showInfo('USER.TOAST.PASSWORDLESSREMOVED', true);
            this.getPasswordless();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }
}
