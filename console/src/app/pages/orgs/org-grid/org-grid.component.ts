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
    constructor(
        public authService: AuthService,
        private userService: AuthUserService,
        private toast: ToastService,
        private router: Router,
    ) {
        this.getData(10, 0);

        this.authService.GetActiveOrg().then(org => this.activeOrg = org);
    }

    private getData(limit: number, offset: number): void {
        this.userService.SearchMyProjectOrgs(limit, offset).then(res => {
            this.orgList = res.toObject().resultList;
            console.log(this.orgList);
        }).catch(error => {
            console.error(error);
            this.toast.showError(error.message);
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
