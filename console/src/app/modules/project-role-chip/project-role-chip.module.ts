import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';

import { ProjectRoleChipComponent } from './project-role-chip.component';

@NgModule({
  declarations: [ProjectRoleChipComponent],
  imports: [CommonModule, MatIconModule, FormsModule, MatButtonModule],
  exports: [ProjectRoleChipComponent],
})
export class ProjectRoleChipModule {}
