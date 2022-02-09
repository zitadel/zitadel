import { Component, Input, OnInit } from '@angular/core';

import { hashCode } from '../avatar/avatar.component';

declare const tinycolor: any;

export class CnslProjectRole {
  private roleName: string = '';
  constructor(roleName: string) {
    this.roleName = roleName;
  }

  public get color(): { hex: string; isLight: boolean; lighter: string; darker: string } {
    const hex = getRoleColor(this.roleName);

    const c = tinycolor(hex);
    return {
      hex: c.toHexString(),
      isLight: c.isLight(),
      lighter: c.lighten(20),
      darker: c.darken(20),
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
  public role!: CnslProjectRole;
  constructor() {}

  public ngOnInit(): void {
    this.role = new CnslProjectRole(this.roleName);
  }
}

export function getRoleColor(roleName: string): string {
  const colors = [
    '#ef4444',
    '#f97316',
    '#f59e0b',
    '#eab308',
    '#84cc16',
    '#22c55e',
    '#10b981',
    '#14b8a6',
    '#06b6d4',
    '#0ea5e9',
    '#3b82f6',
    '#6366f1',
    '#8b5cf6',
    '#a855f7',
    '#d946ef',
    '#ec4899',
    '#f43f5e',
  ];

  let hash = 0;
  if (roleName.length === 0) {
    return colors[hash];
  }

  hash = hashCode(roleName);
  return colors[hash % colors.length];
}
