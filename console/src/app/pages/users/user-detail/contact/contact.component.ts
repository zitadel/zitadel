import { Component, EventEmitter, Input, Output } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Human, UserState } from 'src/app/proto/generated/zitadel/user_pb';

import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { CodeDialogComponent } from '../auth-user-detail/code-dialog/code-dialog.component';
import { EditDialogType } from '../auth-user-detail/edit-dialog/edit-dialog.component';

@Component({
  selector: 'cnsl-contact',
  templateUrl: './contact.component.html',
  styleUrls: ['./contact.component.scss'],
})
export class ContactComponent {
  @Input() disablePhoneCode: boolean = false;
  @Input() canWrite: boolean | null = false;
  @Input() human?: Human.AsObject;
  @Input() username: string = '';
  @Input() state!: UserState;
  @Output() editType: EventEmitter<EditDialogType> = new EventEmitter<EditDialogType>();
  @Output() resendEmailVerification: EventEmitter<void> = new EventEmitter<void>();
  @Output() resendPhoneVerification: EventEmitter<void> = new EventEmitter<void>();
  @Output() enteredPhoneCode: EventEmitter<string> = new EventEmitter<string>();
  @Output() deletedPhone: EventEmitter<void> = new EventEmitter<void>();
  public UserState: any = UserState;

  public EditDialogType: any = EditDialogType;
  constructor(
    private dialog: MatDialog,
    private authService: GrpcAuthService,
  ) {}

  async emitDeletePhone(): Promise<void> {
    const { resultList } = await this.authService.listMyMultiFactors();
    const hasSMSOTP = !!resultList.find((mfa) => mfa.otpSms);

    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'USER.LOGINMETHODS.PHONE.DELETETITLE',
        descriptionKey: 'USER.LOGINMETHODS.PHONE.DELETEDESC',
        warnSectionKey: hasSMSOTP ? 'USER.LOGINMETHODS.PHONE.OTPSMSREMOVALWARNING' : '',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.deletedPhone.emit();
      }
    });
  }

  emitEmailVerification(): void {
    this.resendEmailVerification.emit();
  }

  emitPhoneVerification(): void {
    this.resendPhoneVerification.emit();
  }

  public enterCode(): void {
    if (this.human) {
      const dialogRef = this.dialog.open(CodeDialogComponent, {
        data: {
          number: this.human.phone,
        },
        width: '400px',
      });

      dialogRef.afterClosed().subscribe((code) => {
        if (code) {
          this.enteredPhoneCode.emit(code);
        }
      });
    }
  }

  public openEditDialog(type: EditDialogType): void {
    this.editType.emit(type);
  }
}
