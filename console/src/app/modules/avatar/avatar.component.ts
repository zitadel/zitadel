import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-avatar',
  templateUrl: './avatar.component.html',
  styleUrls: ['./avatar.component.scss'],
})
export class AvatarComponent implements OnInit {
  @Input() name: string = '';
  @Input() credentials: string = '';
  @Input() size: number = 24;
  @Input() fontSize: number = 14;
  @Input() active: boolean = false;
  @Input() color: string = '';
  @Input() forColor: string = '';
  @Input() avatarUrl: string = '';
  constructor() { }

  ngOnInit(): void {
    if (!this.credentials && this.forColor) {
      this.credentials = this.getInitials(this.forColor);
      if (!this.color) {
        this.color = this.getColor(this.forColor || '');
      }
    }

    if (this.size > 50) {
      this.fontSize = 32;
    }
  }

  getInitials(fromName: string): string {
    const username = fromName.split('@')[0];
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

  getColor(userName: string): string {
    const colors = [
      'linear-gradient(40deg, #B44D51 30%, rgb(241,138,138))',
      'linear-gradient(40deg, #B75073 30%, rgb(234,96,143))',
      'linear-gradient(40deg, #84498E 30%, rgb(214,116,230))',
      'linear-gradient(40deg, #705998 30%, rgb(163,131,220))',
      'linear-gradient(40deg, #5C6598 30%, rgb(135,148,222))',
      'linear-gradient(40deg, #7F90D3 30%, rgb(181,196,247))',
      'linear-gradient(40deg, #3E93B9 30%, rgb(150,215,245))',
      'linear-gradient(40deg, #3494A0 30%, rgb(71,205,222))',
      'linear-gradient(40deg, #25716A 30%, rgb(58,185,173))',
      'linear-gradient(40deg, #427E41 30%, rgb(97,185,96))',
      'linear-gradient(40deg, #89A568 30%, rgb(176,212,133))',
      'linear-gradient(40deg, #90924D 30%, rgb(187,189,98))',
      'linear-gradient(40deg, #E2B032 30%, rgb(245,203,99))',
      'linear-gradient(40deg, #C97358 30%, rgb(245,148,118))',
      'linear-gradient(40deg, #6D5B54 30%, rgb(152,121,108))',
      'linear-gradient(40deg, #6B7980 30%, rgb(134,163,177))',
    ];

    let hash = 0;
    if (userName.length === 0) {
      return colors[hash];
    }

    hash = this.hashCode(userName);
    return colors[hash % colors.length];
  }

  // tslint:disable
  private hashCode(str: string, seed: number = 0): number {
    let h1 = 0xdeadbeef ^ seed, h2 = 0x41c6ce57 ^ seed;
    for (let i = 0, ch; i < str.length; i++) {
      ch = str.charCodeAt(i);
      h1 = Math.imul(h1 ^ ch, 2654435761);
      h2 = Math.imul(h2 ^ ch, 1597334677);
    }
    h1 = Math.imul(h1 ^ (h1 >>> 16), 2246822507) ^ Math.imul(h2 ^ (h2 >>> 13), 3266489909);
    h2 = Math.imul(h2 ^ (h2 >>> 16), 2246822507) ^ Math.imul(h1 ^ (h1 >>> 13), 3266489909);
    return 4294967296 * (2097151 & h2) + (h1 >>> 0);
  }
  // tslint:enable
}
