import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTabsModule } from '@angular/material/tabs';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { MemberCreateDialogModule } from 'src/app/modules/add-member-dialog/member-create-dialog.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ContributorsModule } from 'src/app/modules/contributors/contributors.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { PolicyGridModule } from 'src/app/modules/policy-grid/policy-grid.module';
import { SharedModule } from 'src/app/modules/shared/shared.module';
import { WarnDialogModule } from 'src/app/modules/warn-dialog/warn-dialog.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { ChangesModule } from '../../modules/changes/changes.module';
import { AddDomainDialogModule } from './org-detail/add-domain-dialog/add-domain-dialog.module';
import { DomainVerificationComponent } from './org-detail/domain-verification/domain-verification.component';
import { OrgDetailComponent } from './org-detail/org-detail.component';
import { OrgsRoutingModule } from './orgs-routing.module';

@NgModule({
    declarations: [OrgDetailComponent, DomainVerificationComponent],
    imports: [
        CommonModule,
        HasRolePipeModule,
        OrgsRoutingModule,
        FormsModule,
        HasRoleModule,
        InputModule,
        MatButtonModule,
        MatDialogModule,
        CardModule,
        MatIconModule,
        ReactiveFormsModule,
        MatButtonToggleModule,
        MetaLayoutModule,
        MatTabsModule,
        MatTooltipModule,
        WarnDialogModule,
        MemberCreateDialogModule,
        MatMenuModule,
        ChangesModule,
        MatProgressSpinnerModule,
        AddDomainDialogModule,
        TranslateModule,
        SharedModule,
        ContributorsModule,
        CopyToClipboardModule,
        PolicyGridModule,
    ],
})
export class OrgsModule { }
