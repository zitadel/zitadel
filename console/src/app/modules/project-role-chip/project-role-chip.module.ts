import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';

import { ProjectRoleChipComponent } from './project-role-chip.component';

@NgModule({
  declarations: [ProjectRoleChipComponent],
  imports: [CommonModule, MatIconModule, FormsModule, MatButtonModule],
  exports: [ProjectRoleChipComponent],
})
export class ProjectRoleChipModule {}
