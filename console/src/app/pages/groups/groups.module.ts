import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { GroupDetailModule } from './group-detail/group-detail.module';
import { GroupListModule } from './group-list/group-list.module';
import { GroupsRoutingModule } from './groups-routing.module';

@NgModule({
  declarations: [],
  imports: [GroupsRoutingModule, GroupListModule, GroupDetailModule, ChangesModule, CommonModule],
})
export default class GroupsModule {}
