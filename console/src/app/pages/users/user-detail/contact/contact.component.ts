import { Component, EventEmitter, Input, Output } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Human, UserState } from 'src/app/proto/generated/zitadel/user_pb';

import { CodeDialogComponent } from '../auth-user-detail/code-dialog/code-dialog.component';
import { EditDialogType } from '../user-detail/user-detail.component';

@Component({
    selector: 'app-contact',
    templateUrl: './contact.component.html',
    styleUrls: ['./contact.component.scss'],
})
export class ContactComponent {
    @Input() disablePhoneCode: boolean = false;
    @Input() canWrite: boolean = false;
    @Input() human!: Human.AsObject;
    @Input() state!: UserState;
    @Output() editType: EventEmitter<EditDialogType> = new EventEmitter();
    @Output() resendEmailVerification: EventEmitter<void> = new EventEmitter();
    @Output() resendPhoneVerification: EventEmitter<void> = new EventEmitter();
    @Output() enteredPhoneCode: EventEmitter<string> = new EventEmitter();
    @Output() deletedPhone: EventEmitter<void> = new EventEmitter();
    @Input() public userStateEnum: any;

    public EditDialogType: any = EditDialogType;
    constructor(private dialog: MatDialog) { }

    emitDeletePhone(): void {
        const dialogRef = this.dialog.open(WarnDialogComponent, {
            data: {
                confirmKey: 'ACTIONS.DELETE',
                cancelKey: 'ACTIONS.CANCEL',
                titleKey: 'USER.LOGINMETHODS.PHONE.DELETETITLE',
                descriptionKey: 'USER.LOGINMETHODS.PHONE.DELETEDESC',
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
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

            dialogRef.afterClosed().subscribe(code => {
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
