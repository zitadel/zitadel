import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';

import { CardModule } from '../card/card.module';
import { CreateLayoutModule } from '../create-layout/create-layout.module';
import { InfoSectionModule } from '../info-section/info-section.module';
import { ProviderOptionsModule } from '../provider-options/provider-options.module';
import { StringListModule } from '../string-list/string-list.module';
import { LDAPAttributesComponent } from './ldap-attributes/ldap-attributes.component';
import { ProviderAppleComponent } from './provider-apple/provider-apple.component';
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
import { ProviderSamlSpComponent } from './provider-saml-sp/provider-saml-sp.component';
import { CopyRowComponent } from '../../components/copy-row/copy-row.component';
import { ProviderNextComponent } from './provider-next/provider-next.component';
import { ProviderNextService } from './provider-next/provider-next.service';

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
    ProviderAppleComponent,
    ProviderSamlSpComponent,
    ProviderNextComponent,
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
    MatProgressSpinnerModule,
    CopyRowComponent,
  ],
  providers: [ProviderNextService],
})
export default class ProvidersModule {}
