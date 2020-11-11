import { Component, EventEmitter, Input, Output } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { HumanView as AuthHumanView, UserState as AuthUserState } from 'src/app/proto/generated/auth_pb';
import { HumanView as MgmtHumanView, UserState as MgmtUserState } from 'src/app/proto/generated/management_pb';

import { CodeDialogComponent } from '../auth-user-detail/code-dialog/code-dialog.component';

@Component({
    selector: 'app-contact',
    templateUrl: './contact.component.html',
    styleUrls: ['./contact.component.scss'],
})
export class ContactComponent {
    @Input() disablePhoneCode: boolean = false;
    @Input() canWrite: boolean = false;
    @Input() human!: AuthHumanView.AsObject | MgmtHumanView.AsObject;
    @Input() state!: AuthUserState | MgmtUserState;
    @Output() savedPhone: EventEmitter<string> = new EventEmitter();
    @Output() savedEmail: EventEmitter<string> = new EventEmitter();
    @Output() resendEmailVerification: EventEmitter<void> = new EventEmitter();
    @Output() resendPhoneVerification: EventEmitter<void> = new EventEmitter();
    @Output() enteredPhoneCode: EventEmitter<string> = new EventEmitter();
    @Output() deletedPhone: EventEmitter<void> = new EventEmitter();
    @Input() public userStateEnum: any;

    public emailEditState: boolean = false;
    public phoneEditState: boolean = false;
    constructor(private dialog: MatDialog) { }

    savePhone(): void {
        this.phoneEditState = false;
        this.savedPhone.emit(this.human.phone);
    }

    emitDeletePhone(): void {
        this.phoneEditState = false;
        this.deletedPhone.emit();
    }

    saveEmail(): void {
        this.emailEditState = false;
        this.savedEmail.emit(this.human.email);
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
}
