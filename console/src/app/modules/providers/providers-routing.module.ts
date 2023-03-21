import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProviderType } from 'src/app/proto/generated/zitadel/idp_pb';
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

const typeMap = {
  [ProviderType.PROVIDER_TYPE_AZURE_AD]: { path: 'azure', component: ProviderAzureADComponent },
  [ProviderType.PROVIDER_TYPE_GITHUB]: { path: 'github', component: ProviderGithubComponent },
  [ProviderType.PROVIDER_TYPE_GITHUB_ES]: { path: 'github-es', component: ProviderGithubESComponent },
  [ProviderType.PROVIDER_TYPE_GITLAB]: { path: 'gitlab', component: ProviderGitlabComponent },
  [ProviderType.PROVIDER_TYPE_GITLAB_SELF_HOSTED]: {
    path: 'gitlab-self-hosted',
    component: ProviderGitlabSelfHostedComponent,
  },
  [ProviderType.PROVIDER_TYPE_GOOGLE]: { path: 'google', component: ProviderGoogleComponent },
  [ProviderType.PROVIDER_TYPE_JWT]: { path: 'jwt', component: ProviderJWTComponent },
  [ProviderType.PROVIDER_TYPE_OAUTH]: { path: 'oauth', component: ProviderOAuthComponent },
  [ProviderType.PROVIDER_TYPE_OIDC]: { path: 'oidc', component: ProviderOIDCComponent },
  [ProviderType.PROVIDER_TYPE_LDAP]: { path: 'ldap', component: ProviderLDAPComponent },
};

const routes: Routes = Object.entries(typeMap).map(([key, value]) => {
  return {
    path: value.path,
    children: [
      {
        path: 'create',
        component: value.component,
      },
      {
        path: ':id',
        component: value.component,
      },
    ],
  };
});

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProvidersRoutingModule {}
