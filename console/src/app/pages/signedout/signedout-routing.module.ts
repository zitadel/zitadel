import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { SignedoutComponent } from './signedout.component';

const routes: Routes = [
  {
    path: '',
    component: SignedoutComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class SignedoutRoutingModule {}
