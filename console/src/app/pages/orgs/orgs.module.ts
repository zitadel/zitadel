import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatMenuModule } from '@angular/material/menu';
import { MatTabsModule } from '@angular/material/tabs';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { MemberCreateDialogModule } from 'src/app/modules/add-member-dialog/member-create-dialog.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ContributorsModule } from 'src/app/modules/contributors/contributors.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { SharedModule } from 'src/app/modules/shared/shared.module';
import { WarnDialogModule } from 'src/app/modules/warn-dialog/warn-dialog.module';

import { ChangesModule } from '../../modules/changes/changes.module';
import { AddDomainDialogModule } from './org-detail/add-domain-dialog/add-domain-dialog.module';
import { DomainVerificationComponent } from './org-detail/domain-verification/domain-verification.component';
import { OrgDetailComponent } from './org-detail/org-detail.component';
import { OrgGridComponent } from './org-grid/org-grid.component';
import { OrgsRoutingModule } from './orgs-routing.module';
import { PolicyGridComponent } from './policy-grid/policy-grid.component';

@NgModule({
    declarations: [OrgDetailComponent, OrgGridComponent, PolicyGridComponent, DomainVerificationComponent],
    imports: [
        CommonModule,
        OrgsRoutingModule,
        FormsModule,
        HasRoleModule,
        MatFormFieldModule,
        MatInputModule,
        MatButtonModule,
        MatDialogModule,
        CardModule,
        MatIconModule,
        ReactiveFormsModule,
        MatButtonToggleModule,
        MetaLayoutModule,
        MatTabsModule,
        MatTooltipModule,
        MatDialogModule,
        WarnDialogModule,
        MemberCreateDialogModule,
        MatMenuModule,
        ChangesModule,
        AddDomainDialogModule,
        TranslateModule,
        SharedModule,
        ContributorsModule,
        CopyToClipboardModule,
    ],
})
export class OrgsModule { }
