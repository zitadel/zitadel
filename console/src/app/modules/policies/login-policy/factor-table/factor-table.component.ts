import { Component, EventEmitter, Input, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import {
  AddMultiFactorToLoginPolicyRequest as AdminAddMultiFactorToLoginPolicyRequest,
  AddSecondFactorToLoginPolicyRequest as AdminAddSecondFactorToLoginPolicyRequest,
  RemoveMultiFactorFromLoginPolicyRequest as AdminRemoveMultiFactorFromLoginPolicyRequest,
  RemoveSecondFactorFromLoginPolicyRequest as AdminRemoveSecondFactorFromLoginPolicyRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
  AddMultiFactorToLoginPolicyRequest as MgmtAddMultiFactorToLoginPolicyRequest,
  AddSecondFactorToLoginPolicyRequest as MgmtAddSecondFactorToLoginPolicyRequest,
  RemoveMultiFactorFromLoginPolicyRequest as MgmtRemoveMultiFactorFromLoginPolicyRequest,
  RemoveSecondFactorFromLoginPolicyRequest as MgmtRemoveSecondFactorFromLoginPolicyRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { MultiFactorType, SecondFactorType } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { DialogAddTypeComponent } from './dialog-add-type/dialog-add-type.component';

export enum LoginMethodComponentType {
  MultiFactor = 1,
  SecondFactor = 2,
}

@Component({
  selector: 'cnsl-factor-table',
  templateUrl: './factor-table.component.html',
  styleUrls: ['./factor-table.component.scss'],
})
export class FactorTableComponent {
  public LoginMethodComponentType: any = LoginMethodComponentType;
  @Input() componentType!: LoginMethodComponentType;
  @Input() public serviceType!: PolicyComponentServiceType;
  @Input() service!: AdminService | ManagementService;
  @Input() disabled: boolean = false;
  @Input() list: Array<MultiFactorType | SecondFactorType> = [];
  @Output() typeRemoved: EventEmitter<Promise<any>> = new EventEmitter();
  @Output() typeAdded: EventEmitter<Promise<any>> = new EventEmitter();

  @ViewChild(MatPaginator) public paginator!: MatPaginator;

  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();

  public PolicyComponentServiceType: any = PolicyComponentServiceType;

  constructor(
    public translate: TranslateService,
    private toast: ToastService,
    private dialog: MatDialog,
  ) {}

  public removeMfa(type: MultiFactorType | SecondFactorType): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'MFA.DELETE.TITLE',
        descriptionKey: 'MFA.DELETE.DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        let request;

        if (this.serviceType === PolicyComponentServiceType.MGMT) {
          if (this.componentType === LoginMethodComponentType.MultiFactor) {
            const req = new MgmtRemoveMultiFactorFromLoginPolicyRequest();
            req.setType(type as MultiFactorType);
            request = (this.service as ManagementService).removeMultiFactorFromLoginPolicy(req);
          } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
            const req = new MgmtRemoveSecondFactorFromLoginPolicyRequest();
            req.setType(type as SecondFactorType);
            request = (this.service as ManagementService).removeSecondFactorFromLoginPolicy(req);
          }
        } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
          if (this.componentType === LoginMethodComponentType.MultiFactor) {
            const req = new AdminRemoveMultiFactorFromLoginPolicyRequest();
            req.setType(type as MultiFactorType);
            request = (this.service as AdminService).removeMultiFactorFromLoginPolicy(req);
          } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
            const req = new AdminRemoveSecondFactorFromLoginPolicyRequest();
            req.setType(type as SecondFactorType);
            request = (this.service as AdminService).removeSecondFactorFromLoginPolicy(req);
          }
        }

        if (request) {
          this.typeRemoved.emit(request);
        }
      }
    });
  }

  public addMfa(): void {
    const dialogRef = this.dialog.open(DialogAddTypeComponent, {
      data: {
        title: 'MFA.CREATE.TITLE',
        desc: 'MFA.CREATE.DESCRIPTION',
        componentType: this.componentType,
        types: this.availableSelection,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((mfaType: MultiFactorType | SecondFactorType) => {
      if (mfaType) {
        let request;

        if (this.serviceType === PolicyComponentServiceType.MGMT) {
          if (this.componentType === LoginMethodComponentType.MultiFactor) {
            const req = new MgmtAddMultiFactorToLoginPolicyRequest();
            req.setType(mfaType as MultiFactorType);
            request = (this.service as ManagementService).addMultiFactorToLoginPolicy(req);
          } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
            const req = new MgmtAddSecondFactorToLoginPolicyRequest();
            req.setType(mfaType as SecondFactorType);
            request = (this.service as ManagementService).addSecondFactorToLoginPolicy(req);
          }
        } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
          if (this.componentType === LoginMethodComponentType.MultiFactor) {
            const req = new AdminAddMultiFactorToLoginPolicyRequest();
            req.setType(mfaType as MultiFactorType);
            request = (this.service as AdminService).addMultiFactorToLoginPolicy(req);
          } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
            const req = new AdminAddSecondFactorToLoginPolicyRequest();
            req.setType(mfaType as SecondFactorType);
            request = (this.service as AdminService).addSecondFactorToLoginPolicy(req);
          }
        }

        if (request) {
          this.typeAdded.emit(request);
        }
      }
    });
  }

  public get availableSelection(): Array<MultiFactorType | SecondFactorType> {
    const allTypes: MultiFactorType[] | SecondFactorType[] =
      this.componentType === LoginMethodComponentType.MultiFactor
        ? [MultiFactorType.MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION]
        : this.componentType === LoginMethodComponentType.SecondFactor
        ? [
            SecondFactorType.SECOND_FACTOR_TYPE_U2F,
            SecondFactorType.SECOND_FACTOR_TYPE_OTP,
            SecondFactorType.SECOND_FACTOR_TYPE_OTP_SMS,
            SecondFactorType.SECOND_FACTOR_TYPE_OTP_EMAIL,
          ]
        : [];

    const filtered = (allTypes as Array<MultiFactorType | SecondFactorType>).filter((type) => !this.list.includes(type));

    return filtered;
  }
}
