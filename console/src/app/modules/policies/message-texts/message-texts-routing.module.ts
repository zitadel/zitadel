import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { MessageTextsComponent } from './message-texts.component';

const routes: Routes = [
  {
    path: '',
    component: MessageTextsComponent,
    data: {
      animation: 'DetailPage',
    },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class MessageTextsRoutingModule { }
