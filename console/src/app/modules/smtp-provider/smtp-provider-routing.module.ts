import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { SMTPProviderComponent } from './smtp-provider.component';

const routes: Routes = [{ path: ':provider', component: SMTPProviderComponent }];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class SMTPProvidersRoutingModule {}
