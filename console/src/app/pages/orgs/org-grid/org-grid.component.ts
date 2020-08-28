import { SelectionModel } from '@angular/cdk/collections';
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { Observable, of } from 'rxjs';
import { switchMap, take } from 'rxjs/operators';
import { Org } from 'src/app/proto/generated/auth_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-org-grid',
    templateUrl: './org-grid.component.html',
    styleUrls: ['./org-grid.component.scss'],
})
export class OrgGridComponent {
    public activeOrg!: Org.AsObject;
    public orgList: Org.AsObject[] = [];

    public selection: SelectionModel<Org.AsObject> = new SelectionModel<Org.AsObject>(true, []);
    public selectedIndex: number = -1;
    public loading: boolean = false;

    public notPinned: Array<Org.AsObject> = [];

    constructor(
        private userService: GrpcAuthService,
        private toast: ToastService,
        private router: Router,
    ) {
        this.loading = true;
        this.getData(10, 0);

        this.userService.GetActiveOrg().then(org => this.activeOrg = org);

        this.selection.changed.subscribe(selection => {
            this.setPrefixedItem('pinned-orgs', JSON.stringify(
                this.selection.selected.map(item => item.id),
            )).pipe(take(1)).subscribe(() => {
                selection.added.forEach(element => {
                    const index = this.notPinned.findIndex(item => item.id === element.id);
                    this.notPinned.splice(index, 1);
                });

                this.notPinned.push(...selection.removed);
            });
        });
    }

    public reorganizeItems(): void {
        this.getPrefixedItem('pinned-orgs').pipe(take(1)).subscribe(storageEntry => {
            if (storageEntry) {
                const array: string[] = JSON.parse(storageEntry);
                const toSelect: Org.AsObject[] = this.orgList.filter((item, index) => {
                    if (array.includes(item.id)) {
                        // this.notPinned.splice(index, 1);
                        return true;
                    }
                });
                this.selection.select(...toSelect);

                const toNotPinned: Org.AsObject[] = this.orgList.filter((item, index) => {
                    if (!array.includes(item.id)) {
                        return true;
                    }
                });
                this.notPinned = toNotPinned;
            }
        });
    }

    private getPrefixedItem(key: string): Observable<string | null> {
        return this.userService.user.pipe(
            take(1),
            switchMap(user => {
                return of(localStorage.getItem(`${user.id}:${key}`));
            }),
        );
    }

    private setPrefixedItem(key: string, value: any): Observable<void> {
        return this.userService.user.pipe(
            take(1),
            switchMap(user => {
                return of(localStorage.setItem(`${user.id}:${key}`, value));
            }),
        );
    }

    private getData(limit: number, offset: number): void {
        this.userService.SearchMyProjectOrgs(limit, offset).then(res => {
            this.orgList = res.toObject().resultList;

            this.notPinned = Object.assign([], this.orgList);
            this.reorganizeItems();
            this.loading = false;
        }).catch(error => {
            this.toast.showError(error);
            this.loading = false;
        });
    }

    public selectOrg(item: Org.AsObject, event?: any): void {
        if (event && !event.target.classList.contains('mat-icon')) {
            this.userService.setActiveOrg(item);
            this.routeToOrg(item);
        }
    }

    public routeToOrg(item: Org.AsObject): void {
        this.router.navigate(['/orgs', item.id]);
    }
}
