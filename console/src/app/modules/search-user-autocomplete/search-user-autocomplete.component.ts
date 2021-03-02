import { COMMA, ENTER } from '@angular/cdk/keycodes';
import {
    AfterContentChecked,
    ChangeDetectorRef,
    Component,
    ElementRef,
    EventEmitter,
    Input,
    OnInit,
    Output,
    ViewChild,
} from '@angular/core';
import { FormControl } from '@angular/forms';
import { MatAutocomplete, MatAutocompleteSelectedEvent } from '@angular/material/autocomplete';
import { MatChipInputEvent } from '@angular/material/chips';
import { from, of, Subject } from 'rxjs';
import { debounceTime, switchMap, takeUntil, tap } from 'rxjs/operators';
import { SearchMethod, UserSearchKey, UserSearchQuery, UserView } from 'src/app/proto/generated/management_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

export enum UserTarget {
    SELF = 'self',
    EXTERNAL = 'external',
}

@Component({
    selector: 'app-search-user-autocomplete',
    templateUrl: './search-user-autocomplete.component.html',
    styleUrls: ['./search-user-autocomplete.component.scss'],
})
export class SearchUserAutocompleteComponent implements OnInit, AfterContentChecked {
    public selectable: boolean = true;
    public removable: boolean = true;
    public addOnBlur: boolean = true;
    public separatorKeysCodes: number[] = [ENTER, COMMA];

    public myControl: FormControl = new FormControl();
    public globalLoginNameControl: FormControl = new FormControl();

    public loginNames: string[] = [];
    @Input() public users: Array<User.AsObject> = [];
    public filteredUsers: Array<User.AsObject> = [];
    public isLoading: boolean = false;
    @Input() public target: UserTarget = UserTarget.SELF;
    public hint: string = '';
    public UserTarget: any = UserTarget;
    @ViewChild('usernameInput') public usernameInput!: ElementRef<HTMLInputElement>;
    @ViewChild('auto') public matAutocomplete!: MatAutocomplete;
    @Output() public selectionChanged: EventEmitter<User.AsObject | User.AsObject[]> = new EventEmitter();
    @Input() public singleOutput: boolean = false;

    private unsubscribed$: Subject<void> = new Subject();
    constructor(private userService: ManagementService, private toast: ToastService, private cdref: ChangeDetectorRef) { }

    public ngOnInit(): void {
        if (this.target === UserTarget.EXTERNAL) {
            this.filteredUsers = [];
            this.unsubscribed$.next(); // clear old subscription
        } else if (this.target === UserTarget.SELF) {
            this.getFilteredResults(); // new subscription
        }
    }

    public ngAfterContentChecked(): void {
        this.cdref.detectChanges();
    }

    private getFilteredResults(): void {
        this.myControl.valueChanges.pipe(debounceTime(200),
            takeUntil(this.unsubscribed$),
            tap(() => this.isLoading = true),
            switchMap(value => {
                const query = new UserSearchQuery();
                query.setKey(UserSearchKey.USERSEARCHKEY_USER_NAME);
                query.setValue(value);
                query.setMethod(SearchMethod.SEARCHMETHOD_CONTAINS_IGNORE_CASE);
                if (this.target === UserTarget.SELF) {
                    return from(this.userService.listUsers(10, 0, [query]));
                } else {
                    return of(); // from(this.userService.GetUserByEmailGlobal(value));
                }
            }),
        ).subscribe((userresp: any) => {
            this.isLoading = false;
            if (this.target === UserTarget.SELF && userresp) {
                this.filteredUsers = userresp.toObject().resultList;
            }
        });
    }

    public displayFn(user?: UserView.AsObject): string | undefined {
        return user ? `${user.preferredLoginName}` : undefined;
    }

    public add(event: MatChipInputEvent): void {
        if (!this.matAutocomplete.isOpen) {
            const input = event.input;
            const value = event.value;

            if ((value || '').trim()) {
                const index = this.filteredUsers.findIndex((user) => {
                    if (user.preferredLoginName) {
                        return user.preferredLoginName === value;
                    }
                });
                if (index > -1) {
                    if (this.users && this.users.length > 0) {
                        this.users.push(this.filteredUsers[index]);
                    } else {
                        this.users = [this.filteredUsers[index]];
                    }
                }
            }

            if (input) {
                input.value = '';
            }
        }
    }

    public remove(user: UserView.AsObject): void {
        const index = this.users.indexOf(user);

        if (index >= 0) {
            this.users.splice(index, 1);
            this.selectionChanged.emit(this.users);
        }
    }

    public selected(event: MatAutocompleteSelectedEvent): void {
        const index = this.filteredUsers.findIndex((user) => user === event.option.value);
        if (index !== -1) {
            if (this.singleOutput) {
                this.selectionChanged.emit(this.filteredUsers[index]);
            } else {
                if (this.users && this.users.length > 0) {
                    this.users.push(this.filteredUsers[index]);
                } else {
                    this.users = [this.filteredUsers[index]];
                }

                if (this.singleOutput) {
                    this.selectionChanged.emit(this.users[0]);
                } else {
                    this.selectionChanged.emit(this.users);
                }
                this.usernameInput.nativeElement.value = '';
                this.myControl.setValue(null);
            }
        }
    }

    public changeTarget(): void {
        if (this.target === UserTarget.SELF) {
            this.target = UserTarget.EXTERNAL;
            this.filteredUsers = [];
            this.unsubscribed$.next(); // clear old subscription
        } else if (this.target === UserTarget.EXTERNAL) {
            this.target = UserTarget.SELF;
            this.getFilteredResults(); // new subscription
        }
    }

    public getGlobalUser(): void {
        this.userService.getUserByLoginNameGlobal(this.globalLoginNameControl.value).then(resp => {
            if (this.singleOutput && resp.user) {
                this.users = [resp.user];
                this.selectionChanged.emit(this.users[0]);
            } else if (resp.user) {
                this.users.push(resp.user);
                this.selectionChanged.emit(this.users);
            }
        }).catch(error => {
            this.toast.showError(error);
        });
    }
}
