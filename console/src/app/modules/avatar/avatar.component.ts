import { Component, Input, OnInit } from '@angular/core';

interface Color {
  200: string;
  300: string;
  500: string;
  900: string;
}

@Component({
  selector: 'cnsl-avatar',
  templateUrl: './avatar.component.html',
  styleUrls: ['./avatar.component.scss'],
})
export class AvatarComponent implements OnInit {
  @Input() name: string = '';
  @Input() credentials: string = '';
  @Input() size: number = 32;
  @Input() fontSize: number = 14;
  @Input() fontWeight: number = 600;
  @Input() active: boolean = false;
  @Input() forColor: string = '';
  @Input() avatarUrl: string = '';
  @Input() isMachine: boolean = false;

  constructor() {}

  ngOnInit(): void {
    if (!this.credentials && this.forColor) {
      this.credentials = this.getInitials(this.forColor);
    } else if (!this.credentials && this.name) {
      this.credentials = this.getInitials(this.name);
    }

    if (this.size > 50) {
      this.fontSize = 32;
      this.fontWeight = 500;
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

  public get color(): Color {
    const toGen = this.forColor || this.name || '';
    console.log(toGen);
    return getColorHash(toGen);
  }
}

export function hashCode(str: string, seed: number = 0): number {
  let h1 = 0xdeadbeef ^ seed,
    h2 = 0x41c6ce57 ^ seed;
  for (let i = 0, ch; i < str.length; i++) {
    ch = str.charCodeAt(i);
    h1 = Math.imul(h1 ^ ch, 2654435761);
    h2 = Math.imul(h2 ^ ch, 1597334677);
  }
  h1 = Math.imul(h1 ^ (h1 >>> 16), 2246822507) ^ Math.imul(h2 ^ (h2 >>> 13), 3266489909);
  h2 = Math.imul(h2 ^ (h2 >>> 16), 2246822507) ^ Math.imul(h1 ^ (h1 >>> 13), 3266489909);
  return 4294967296 * (2097151 & h2) + (h1 >>> 0);
}

function getColorHash(value: string): Color {
  const colors = [
    { 500: '#ef4444', 200: '#fecaca', 300: '#fca5a5', 900: '#7f1d1d' },
    { 500: '#f97316', 200: '#fed7aa', 300: '#fdba74', 900: '#7c2d12' },
    { 500: '#f59e0b', 200: '#fde68a', 300: '#fcd34d', 900: '#78350f' },
    { 500: '#eab308', 200: '#fef08a', 300: '#fde047', 900: '#713f12' },
    { 500: '#84cc16', 200: '#d9f99d', 300: '#bef264', 900: '#365314' },
    { 500: '#22c55e', 200: '#bbf7d0', 300: '#86efac', 900: '#14532d' },
    { 500: '#10b981', 200: '#a7f3d0', 300: '#6ee7b7', 900: '#064e3b' },
    { 500: '#14b8a6', 200: '#99f6e4', 300: '#5eead4', 900: '#134e4a' },
    { 500: '#06b6d4', 200: '#a5f3fc', 300: '#67e8f9', 900: '#164e63' },
    { 500: '#0ea5e9', 200: '#bae6fd', 300: '#7dd3fc', 900: '#0c4a6e' },
    { 500: '#3b82f6', 200: '#bfdbfe', 300: '#93c5fd', 900: '#1e3a8a' },
    { 500: '#6366f1', 200: '#c7d2fe', 300: '#a5b4fc', 900: '#312e81' },
    { 500: '#8b5cf6', 200: '#ddd6fe', 300: '#c4b5fd', 900: '#4c1d95' },
    { 500: '#a855f7', 200: '#e9d5ff', 300: '#d8b4fe', 900: '#581c87' },
    { 500: '#d946ef', 200: '#f5d0fe', 300: '#f0abfc', 900: '#701a75' },
    { 500: '#ec4899', 200: '#fbcfe8', 300: '#f9a8d4', 900: '#831843' },
    { 500: '#f43f5e', 200: '#fecdd3', 300: '#fda4af', 900: '#881337' },
  ];

  let hash = 0;

  if (value.length === 0) {
    return colors[hash];
  }

  hash = hashCode(value);
  return colors[hash % colors.length];
}

export function getColor(toGen: string): string {
  const color = getColorHash(toGen);
  return color[500];
}
