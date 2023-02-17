import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { UserDetailModule } from './user-detail/user-detail.module';
import { UserListModule } from './user-list/user-list.module';
import { UsersRoutingModule } from './users-routing.module';

@NgModule({
  declarations: [],
  imports: [UsersRoutingModule, UserListModule, UserDetailModule, ChangesModule, CommonModule],
})
export default class UsersModule {}
