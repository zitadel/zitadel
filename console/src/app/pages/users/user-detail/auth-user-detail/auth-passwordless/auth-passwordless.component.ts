import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSort } from '@angular/material/sort';
import { MatTable, MatTableDataSource } from '@angular/material/table';
import { BehaviorSubject, Observable } from 'rxjs';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { AuthFactorState, WebAuthNToken } from 'src/app/proto/generated/zitadel/user_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

import { _base64ToArrayBuffer } from '../../u2f-util';
import { U2FComponentDestination } from '../dialog-u2f/dialog-u2f.component';
import { DialogPasswordlessComponent } from './dialog-passwordless/dialog-passwordless.component';

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
  selector: 'app-auth-passwordless',
  templateUrl: './auth-passwordless.component.html',
  styleUrls: ['./auth-passwordless.component.scss'],
})
export class AuthPasswordlessComponent implements OnInit, OnDestroy {
  public displayedColumns: string[] = ['name', 'state', 'actions'];
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();

  @ViewChild(MatTable) public table!: MatTable<WebAuthNToken.AsObject>;
  @ViewChild(MatSort) public sort!: MatSort;
  public dataSource!: MatTableDataSource<WebAuthNToken.AsObject>;

  public AuthFactorState: any = AuthFactorState;
  public error: string = '';

  constructor(private service: GrpcAuthService,
    private toast: ToastService,
    private dialog: MatDialog) { }

  public ngOnInit(): void {
    this.getPasswordless();
  }

  public ngOnDestroy(): void {
    this.loadingSubject.complete();
  }

  public addPasswordless(): void {
    this.service.addMyPasswordless().then((resp) => {
      if (resp.key) {
        const credOptions: CredentialCreationOptions = JSON.parse(atob(resp.key.publicKey as string));

        if (credOptions.publicKey?.challenge) {
          credOptions.publicKey.challenge = _base64ToArrayBuffer(credOptions.publicKey.challenge as any);
          credOptions.publicKey.user.id = _base64ToArrayBuffer(credOptions.publicKey.user.id as any);
          if (credOptions.publicKey.excludeCredentials) {
            credOptions.publicKey.excludeCredentials.map(cred => {
              cred.id = _base64ToArrayBuffer(cred.id as any);
              return cred;
            });
          }
          const dialogRef = this.dialog.open(DialogPasswordlessComponent, {
            width: '400px',
            data: {
              credOptions,
              type: U2FComponentDestination.PASSWORDLESS,
            },
          });

          dialogRef.afterClosed().subscribe(done => {
            this.getPasswordless();
          });
        }
      }
    }, error => {
      this.toast.showError(error);
    });
  }

  public getPasswordless(): void {
    this.service.listMyPasswordless().then(passwordless => {
      this.dataSource = new MatTableDataSource(passwordless.resultList);
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
        this.service.removeMyPasswordless(id).then(() => {
          this.toast.showInfo('USER.TOAST.PASSWORDLESSREMOVED', true);
          this.getPasswordless();
        }).catch(error => {
          this.toast.showError(error);
        });
      }
    });
  }
}
