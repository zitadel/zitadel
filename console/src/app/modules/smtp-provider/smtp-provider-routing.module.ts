import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { SMTPProviderComponent } from './smtp-provider.component';

const types = [
  { path: 'aws-ses', component: SMTPProviderComponent },
  { path: 'generic', component: SMTPProviderComponent },
  { path: 'google', component: SMTPProviderComponent },
  { path: 'mailgun', component: SMTPProviderComponent },
  { path: 'postmark', component: SMTPProviderComponent },
  { path: 'sendgrid', component: SMTPProviderComponent },
  { path: 'mailjet', component: SMTPProviderComponent },
  { path: 'mailchimp', component: SMTPProviderComponent },
  { path: 'brevo', component: SMTPProviderComponent },
  { path: 'outlook', component: SMTPProviderComponent },
];

const routes: Routes = types.map((value) => {
  return {
    path: value.path,
    children: [
      {
        path: 'create',
        component: value.component,
      },
    ],
  };
});

routes.push({
  path: ':id',
  component: SMTPProviderComponent,
});

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class SMTPProvidersRoutingModule {}
