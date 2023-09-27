import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { SMTPProviderType } from 'src/app/proto/generated/zitadel/settings_pb';
import { SMTPProviderComponent } from './smtp-provider.component';

const typeMap = {
  [SMTPProviderType.SMTP_PROVIDER_TYPE_AMAZONSES]: { path: 'aws-ses', component: SMTPProviderComponent },
  [SMTPProviderType.SMTP_PROVIDER_TYPE_GENERIC]: { path: 'generic', component: SMTPProviderComponent },
  [SMTPProviderType.SMTP_PROVIDER_TYPE_GOOGLE]: { path: 'google', component: SMTPProviderComponent },
  [SMTPProviderType.SMTP_PROVIDER_TYPE_MAILGUN]: { path: 'mailgun', component: SMTPProviderComponent },
  [SMTPProviderType.SMTP_PROVIDER_TYPE_POSTMARK]: { path: 'postmark', component: SMTPProviderComponent },
  [SMTPProviderType.SMTP_PROVIDER_TYPE_SENDGRID]: { path: 'sendgrid', component: SMTPProviderComponent },
  [SMTPProviderType.SMTP_PROVIDER_TYPE_MAILJET]: { path: 'mailjet', component: SMTPProviderComponent },
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
export class SMTPProvidersRoutingModule {}
