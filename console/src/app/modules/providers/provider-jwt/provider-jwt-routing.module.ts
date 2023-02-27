import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ProviderJWTComponent } from './provider-jwt.component';

const routes: Routes = [
  {
    path: '',
    component: ProviderJWTComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProviderJWTCreateRoutingModule {}
