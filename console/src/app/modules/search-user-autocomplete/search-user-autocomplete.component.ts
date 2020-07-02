import { COMMA, ENTER } from '@angular/cdk/keycodes';
import { Component, ElementRef, EventEmitter, Input, Output, ViewChild } from '@angular/core';
import { FormControl } from '@angular/forms';
import { MatAutocomplete, MatAutocompleteSelectedEvent } from '@angular/material/autocomplete';
import { MatChipInputEvent } from '@angular/material/chips';
import { from, of, Subject } from 'rxjs';
import { debounceTime, switchMap, takeUntil, tap } from 'rxjs/operators';
import { SearchMethod, User, UserSearchKey, UserSearchQuery } from 'src/app/proto/generated/management_pb';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';
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
export class SearchUserAutocompleteComponent {
    public selectable: boolean = true;
    public removable: boolean = true;
    public addOnBlur: boolean = true;
    public separatorKeysCodes: number[] = [ENTER, COMMA];

    public myControl: FormControl = new FormControl();
    public globalEmailControl: FormControl = new FormControl();

    public emails: string[] = [];
    public users: Array<User.AsObject> = [];
    public filteredUsers: Array<User.AsObject> = [];
    public isLoading: boolean = false;
    public target: UserTarget = UserTarget.SELF;
    public hint: string = '';
    public UserTarget: any = UserTarget;
    @ViewChild('usernameInput') public usernameInput!: ElementRef<HTMLInputElement>;
    @ViewChild('auto') public matAutocomplete!: MatAutocomplete;
    @Output() public selectionChanged: EventEmitter<User.AsObject | User.AsObject[]> = new EventEmitter();
    @Input() public singleOutput: boolean = false;
    private unsubscribed$: Subject<void> = new Subject();
    constructor(private userService: MgmtUserService, private toast: ToastService) {
        this.getFilteredResults();
    }

    private getFilteredResults(): void {
        this.myControl.valueChanges.pipe(debounceTime(200),
            takeUntil(this.unsubscribed$),
            tap(() => this.isLoading = true),
            switchMap(value => {
                const query = new UserSearchQuery();
                query.setKey(UserSearchKey.USERSEARCHKEY_EMAIL);
                query.setValue(value);
                query.setMethod(SearchMethod.SEARCHMETHOD_CONTAINS_IGNORE_CASE);
                if (this.target === UserTarget.SELF) {
                    return from(this.userService.SearchUsers(10, 0, [query]));
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

    public displayFn(user?: User.AsObject): string | undefined {
        return user ? `${user.email}` : undefined;
    }

    public add(event: MatChipInputEvent): void {
        if (!this.matAutocomplete.isOpen) {
            const input = event.input;
            const value = event.value;

            if ((value || '').trim()) {
                const index = this.filteredUsers.findIndex((user) => {
                    if (user.email) {
                        return user.email === value;
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

    public remove(user: User.AsObject): void {
        const index = this.users.indexOf(user);

        if (index >= 0) {
            this.users.splice(index, 1);
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
                this.selectionChanged.emit(this.users);

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
        this.userService.GetUserByEmailGlobal(this.globalEmailControl.value).then(user => {
            this.users = [user.toObject()];
            this.selectionChanged.emit(this.users);
        }).catch(error => {
            this.toast.showError(error.message);
        });
    }
}
