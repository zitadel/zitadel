import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

import { ProjectRoleChipComponent } from './project-role-chip.component';

@NgModule({
  declarations: [ProjectRoleChipComponent],
  imports: [CommonModule, MatIconModule, FormsModule, MatButtonModule],
  exports: [ProjectRoleChipComponent],
})
export class ProjectRoleChipModule {}
