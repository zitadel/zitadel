import { SelectionModel } from '@angular/cdk/collections';
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { Org } from 'src/app/proto/generated/auth_pb';
import { AuthUserService } from 'src/app/services/auth-user.service';
import { AuthService } from 'src/app/services/auth.service';
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
        public authService: AuthService,
        private userService: AuthUserService,
        private toast: ToastService,
        private router: Router,
    ) {
        this.loading = true;
        this.getData(10, 0);

        this.authService.GetActiveOrg().then(org => this.activeOrg = org);

        this.selection.changed.subscribe(selection => {
            this.setPrefixedItem('pinned-orgs', JSON.stringify(
                this.selection.selected.map(item => item.id),
            )).then(() => {
                const filtered = this.notPinned.filter(item => item === selection.added.find(i => i === item));
                filtered.forEach((f, i) => {
                    this.notPinned.splice(i, 1);
                });

                this.notPinned.push(...selection.removed);
            });
        });
    }

    public reorganizeItems(): void {
        this.getPrefixedItem('pinned-orgs').then(storageEntry => {
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

    private async getPrefixedItem(key: string): Promise<string | null> {
        const prefix = (await this.authService.GetActiveOrg()).id;
        return localStorage.getItem(`${prefix}:${key}`);
    }

    private async setPrefixedItem(key: string, value: any): Promise<void> {
        const prefix = (await this.authService.GetActiveOrg()).id;
        return localStorage.setItem(`${prefix}:${key}`, value);
    }

    private getData(limit: number, offset: number): void {
        this.userService.SearchMyProjectOrgs(limit, offset).then(res => {
            this.orgList = res.toObject().resultList;

            this.notPinned = Object.assign([], this.orgList);
            this.reorganizeItems();
            this.loading = false;
        }).catch(error => {
            console.error(error);
            this.toast.showError(error);
            this.loading = false;
        });
    }

    public selectOrg(item: Org.AsObject, event?: any): void {
        if (event && !event.target.classList.contains('mat-icon')) {
            this.authService.setActiveOrg(item);
            this.routeToOrg(item);
        }
    }

    public routeToOrg(item: Org.AsObject): void {
        this.router.navigate(['/orgs', item.id]);
    }
}
