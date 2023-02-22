import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ProviderJWTCreateComponent } from './provider-jwt-create.component';

const routes: Routes = [
  {
    path: '',
    component: ProviderJWTCreateComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProviderJWTCreateRoutingModule {}
