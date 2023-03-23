import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyChipsModule as MatChipsModule } from '@angular/material/legacy-chips';
import { MatLegacyProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';

import { CardModule } from '../card/card.module';
import { CreateLayoutModule } from '../create-layout/create-layout.module';
import { InfoSectionModule } from '../info-section/info-section.module';
import { ProviderOptionsModule } from '../provider-options/provider-options.module';
import { StringListModule } from '../string-list/string-list.module';
import { LDAPAttributesComponent } from './ldap-attributes/ldap-attributes.component';
import { ProviderAzureADComponent } from './provider-azure-ad/provider-azure-ad.component';
import { ProviderGithubESComponent } from './provider-github-es/provider-github-es.component';
import { ProviderGithubComponent } from './provider-github/provider-github.component';
import { ProviderGitlabSelfHostedComponent } from './provider-gitlab-self-hosted/provider-gitlab-self-hosted.component';
import { ProviderGitlabComponent } from './provider-gitlab/provider-gitlab.component';
import { ProviderGoogleComponent } from './provider-google/provider-google.component';
import { ProviderJWTComponent } from './provider-jwt/provider-jwt.component';
import { ProviderLDAPComponent } from './provider-ldap/provider-ldap.component';
import { ProviderOAuthComponent } from './provider-oauth/provider-oauth.component';
import { ProviderOIDCComponent } from './provider-oidc/provider-oidc.component';
import { ProvidersRoutingModule } from './providers-routing.module';

@NgModule({
  declarations: [
    ProviderGoogleComponent,
    ProviderGithubComponent,
    ProviderGithubESComponent,
    ProviderAzureADComponent,
    LDAPAttributesComponent,
    ProviderGitlabSelfHostedComponent,
    ProviderGitlabComponent,
    ProviderGithubESComponent,
    ProviderJWTComponent,
    ProviderOIDCComponent,
    ProviderOAuthComponent,
    ProviderLDAPComponent,
  ],
  imports: [
    ProvidersRoutingModule,
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    CreateLayoutModule,
    StringListModule,
    InfoSectionModule,
    InputModule,
    MatButtonModule,
    MatSelectModule,
    MatIconModule,
    MatChipsModule,
    CardModule,
    MatCheckboxModule,
    MatTooltipModule,
    TranslateModule,
    ProviderOptionsModule,
    MatLegacyProgressSpinnerModule,
  ],
})
export default class ProvidersModule {}
