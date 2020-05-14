import { SelectionModel } from '@angular/cdk/collections';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatButtonToggleChange } from '@angular/material/button-toggle';
import { MatTableDataSource } from '@angular/material/table';
import { ActivatedRoute, Params } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { Org, OrgMember, OrgMemberSearchResponse, OrgState } from 'src/app/proto/generated/management_pb';
import { OrgService } from 'src/app/services/org.service';
import { ToastService } from 'src/app/services/toast.service';


@Component({
    selector: 'app-org-detail',
    templateUrl: './org-detail.component.html',
    styleUrls: ['./org-detail.component.scss'],
})
export class OrgDetailComponent implements OnInit, OnDestroy {
    public orgId: string = '';
    public org!: Org.AsObject;

    public dataSource: MatTableDataSource<OrgMember.AsObject> = new MatTableDataSource<OrgMember.AsObject>();
    public memberResult!: OrgMemberSearchResponse.AsObject;
    public displayedColumns: string[] = ['select', 'firstname', 'lastname', 'username', 'email', 'roles'];
    public selection: SelectionModel<OrgMember.AsObject> = new SelectionModel<OrgMember.AsObject>(true, []);
    public OrgState: any = OrgState;
    public ChangeType: any = ChangeType;

    private subscription: Subscription = new Subscription();

    constructor(
        public translate: TranslateService,
        private route: ActivatedRoute,
        private orgService: OrgService,
        private toast: ToastService,
    ) { }

    public ngOnInit(): void {
        this.subscription = this.route.params.subscribe(params => this.getData(params));
    }

    public ngOnDestroy(): void {
        this.subscription.unsubscribe();
    }

    private async getData({ id }: Params): Promise<void> {
        this.orgId = id;

        this.orgService.GetOrgById(this.orgId).then((org: Org) => {
            this.org = org.toObject();
        }).catch(error => {
            this.toast.showError(error.message);
        });
    }

    public changeState(event: MatButtonToggleChange | any): void {
        if (event.value === OrgState.ORGSTATE_ACTIVE) {
            this.orgService.ReactivateOrg(this.orgId).then(() => {
                this.toast.showInfo('Reactivated Org');
            }).catch(error => {
                this.toast.showError(error.message);
            });
        } else if (event.value === OrgState.ORGSTATE_INACTIVE) {
            this.orgService.DeactivateOrg(this.orgId).then(() => {
                this.toast.showInfo('Deactivated Org');
            }).catch(error => {
                this.toast.showError(error.message);
            });
        }
    }
}
