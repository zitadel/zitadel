import { LiveAnnouncer } from '@angular/cdk/a11y';
import { CommonModule } from '@angular/common';
import { Component, OnDestroy, ViewChild, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialog } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatSort, Sort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { TranslateModule } from '@ngx-translate/core';
import { BehaviorSubject, Observable, Subject, takeUntil } from 'rxjs';
import { CardModule } from 'src/app/modules/card/card.module';
import { DisplayJsonDialogComponent } from 'src/app/modules/display-json-dialog/display-json-dialog.component';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { ListEventsRequest, ListEventsResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { Event } from 'src/app/proto/generated/zitadel/event_pb';
import { FeatureServiceClient } from 'src/app/proto/generated/zitadel/feature/v2beta/Feature_serviceServiceClientPb';
import {
  GetInstanceFeaturesResponse,
  SetInstanceFeaturesRequest,
} from 'src/app/proto/generated/zitadel/feature/v2beta/instance_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { FeatureService } from 'src/app/services/feature.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  imports: [
    CommonModule,
    FormsModule,
    HasRolePipeModule,
    MatIconModule,
    CardModule,
    TranslateModule,
    MatButtonModule,
    MatCheckboxModule,
    InfoSectionModule,
  ],
  standalone: true,
  selector: 'cnsl-features',
  templateUrl: './features.component.html',
  styleUrls: ['./features.component.scss'],
})
export class FeaturesComponent implements OnDestroy {
  private destroy$: Subject<void> = new Subject();

  public _loading: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public featureData: GetInstanceFeaturesResponse.AsObject | undefined = undefined;

  constructor(
    private featureService: FeatureService,
    private breadcrumbService: BreadcrumbService,
    private toast: ToastService,
    private dialog: MatDialog,
  ) {
    const breadcrumbs = [
      new Breadcrumb({
        type: BreadcrumbType.INSTANCE,
        name: 'Instance',
        routerLink: ['/instance'],
      }),
    ];
    this.breadcrumbService.setBreadcrumb(breadcrumbs);

    this.getFeatures(true);
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public openDialog(event: Event): void {
    this.dialog.open(DisplayJsonDialogComponent, {
      data: {
        event: event,
      },
      width: '450px',
    });
  }

  private getFeatures(inheritance: boolean) {
    this.featureService.getInstanceFeatures(inheritance).then((instanceFeaturesResponse) => {
      this.featureData = instanceFeaturesResponse.toObject();
    });
  }

  public saveFeatures(): void {
    if (this.featureData) {
      const req = new SetInstanceFeaturesRequest();
      req.setLoginDefaultOrg(!!this.featureData.loginDefaultOrg?.enabled);
      req.setOidcLegacyIntrospection(!!this.featureData.oidcLegacyIntrospection?.enabled);
      req.setOidcTokenExchange(!!this.featureData.oidcTokenExchange?.enabled);
      req.setOidcTriggerIntrospectionProjections(!!this.featureData.oidcTriggerIntrospectionProjections?.enabled);
      req.setUserSchema(!!this.featureData.userSchema?.enabled);

      this.featureService
        .setInstanceFeatures(req)
        .then(() => {
          this.toast.showInfo('POLICY.TOAST.SET', true);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }
}
