import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { MFAState, MfaType, MultiFactor, UserView } from 'src/app/proto/generated/management_pb';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';


export interface MFAItem {
    name: string;
    verified: boolean;
}

@Component({
    selector: 'app-user-mfa',
    templateUrl: './user-mfa.component.html',
    styleUrls: ['./user-mfa.component.scss'],
})
export class UserMfaComponent implements OnInit, OnDestroy {
    @Input() private user!: UserView.AsObject;
    public mfaSubject: BehaviorSubject<MultiFactor.AsObject[]> = new BehaviorSubject<MultiFactor.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    public MfaType: any = MfaType;
    public MFAState: any = MFAState;

    public error: string = '';
    constructor(private mgmtUserService: MgmtUserService) { }

    public ngOnInit(): void {
        this.getOTP();
    }

    public ngOnDestroy(): void {
        this.mfaSubject.complete();
        this.loadingSubject.complete();
    }

    public getOTP(): void {
        this.mgmtUserService.getUserMfas(this.user.id).then(mfas => {
            this.mfaSubject.next(mfas.toObject().mfasList);
            this.error = '';
        }).catch(error => {
            console.error(error);
            this.error = error.message;
        });
    }
}
