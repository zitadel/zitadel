import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { ThemeService } from 'src/app/services/theme.service';
import { getColorHash } from 'src/app/utils/color';

declare const tinycolor: any;

export class CnslProjectRole {
  private roleName: string = '';
  constructor(roleName: string) {
    this.roleName = roleName;
  }

  public get color(): { hex: string; isLight: boolean; lighter: string; darker: string } {
    const color = getColorHash(this.roleName);
    const hex = color[500];

    const c = tinycolor(hex);
    return {
      hex: c.toHexString(),
      isLight: c.isLight(),
      lighter: color[300],
      darker: color[600],
    };
  }
}

@Component({
  selector: 'cnsl-project-role-chip',
  templateUrl: './project-role-chip.component.html',
  styleUrls: ['./project-role-chip.component.scss'],
})
export class ProjectRoleChipComponent implements OnInit {
  @Input() roleName: string = '';
  @Output() removed: EventEmitter<void> = new EventEmitter();
  @Input() showRemove: boolean = false;
  public role!: CnslProjectRole;
  constructor(public themeService: ThemeService) {}

  public ngOnInit(): void {
    this.role = new CnslProjectRole(this.roleName);
  }

  public emitRemove(): void {
    this.removed.emit();
  }
}
