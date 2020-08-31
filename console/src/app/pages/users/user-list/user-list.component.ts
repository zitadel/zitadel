import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable, Subscription } from 'rxjs';
import { User, UserSearchResponse } from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

export enum UserType {
    HUMAN = 'human',
    MACHINE = 'machine',
}
@Component({
    selector: 'app-user-list',
    templateUrl: './user-list.component.html',
    styleUrls: ['./user-list.component.scss'],
})
export class UserListComponent implements OnDestroy {
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    public dataSource: MatTableDataSource<User.AsObject> = new MatTableDataSource<User.AsObject>();
    public userResult!: UserSearchResponse.AsObject;
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    public displayedColumns: string[] = ['select', 'firstname', 'lastname', 'username', 'email', 'state'];
    public selection: SelectionModel<User.AsObject> = new SelectionModel<User.AsObject>(true, []);
    @Output() public changedSelection: EventEmitter<Array<User.AsObject>> = new EventEmitter();

    private subscription?: Subscription;

    constructor(public translate: TranslateService, private route: ActivatedRoute, private userService: ManagementService,
        private toast: ToastService) {
        this.subscription = this.route.params.subscribe(() => this.getData(10, 0));

        this.selection.changed.subscribe(() => {
            this.changedSelection.emit(this.selection.selected);
        });
    }
}
