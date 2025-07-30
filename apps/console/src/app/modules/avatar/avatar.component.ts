import { Component, Input, OnInit } from '@angular/core';
import { ThemeService } from 'src/app/services/theme.service';
import { Color, getColorHash } from 'src/app/utils/color';

@Component({
  selector: 'cnsl-avatar',
  templateUrl: './avatar.component.html',
  styleUrls: ['./avatar.component.scss'],
})
export class AvatarComponent implements OnInit {
  @Input() name: string = '';
  @Input() size: number = 32;
  @Input() fontSize: number = 14;
  @Input() fontWeight: number = 600;
  @Input() active: boolean = false;
  @Input() forColor: string = '';
  @Input() avatarUrl: string = '';
  @Input() isMachine: boolean = false;

  constructor(public themeService: ThemeService) {}

  ngOnInit(): void {
    if (this.size > 50) {
      this.fontSize = 32;
      this.fontWeight = 500;
    }
  }

  public get credentials(): string {
    const toSplit = this.name ? this.name : this.forColor;

    if (this.name) {
      const split = toSplit.split(' ');
      const initials = split[0].charAt(0) + (split[1] ? split[1].charAt(0) : '');
      return initials;
    } else {
      const username = toSplit.split('@')[0];
      let separator = '_';
      if (username.includes('-')) {
        separator = '-';
      }
      if (username.includes('.')) {
        separator = '.';
      }
      const split = username.split(separator);
      const initials = split[0].charAt(0) + (split[1] ? split[1].charAt(0) : '');
      return initials;
    }
  }

  public get color(): Color {
    const toGen = this.forColor || this.name || '';
    return getColorHash(toGen);
  }

  public errorHandler(event: any) {
    (event.target as HTMLImageElement).style.display = 'none';
  }
}
