import { Component, Input, OnInit } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { UserProfile } from 'src/app/proto/generated/management_pb';


export interface MFAItem {
    name: string;
    verified: boolean;
}

@Component({
    selector: 'app-user-mfa',
    templateUrl: './user-mfa.component.html',
    styleUrls: ['./user-mfa.component.scss'],
})
export class UserMfaComponent implements OnInit {
    @Input() public profile!: UserProfile;

    public mfaSubject: BehaviorSubject<MFAItem[]> = new BehaviorSubject<MFAItem[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    constructor() { }

    ngOnInit(): void {


    }
}
