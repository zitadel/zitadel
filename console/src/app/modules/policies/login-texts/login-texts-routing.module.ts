import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { LoginTextsComponent } from './login-texts.component';

const routes: Routes = [
  {
    path: '',
    component: LoginTextsComponent,
    data: {
      animation: 'DetailPage',
    },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class LoginTextsRoutingModule { }
