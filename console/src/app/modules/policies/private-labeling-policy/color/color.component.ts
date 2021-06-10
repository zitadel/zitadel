import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { ColorEvent } from 'ngx-color';

import { ColorType } from '../private-labeling-policy.component';

@Component({
  selector: 'cnsl-color',
  templateUrl: './color.component.html',
  styleUrls: ['./color.component.scss'],
})
export class ColorComponent implements OnInit {
  public PRIMARY: Array<{ name: string; color: string; }> = [
    { name: 'red', color: '#f44336' },
    { name: 'pink', color: '#e91e63' },
    { name: 'purple', color: '#9c27b0' },
    { name: 'deeppurple', color: '#673ab7' },
    { name: 'indigo', color: '#3f51b5' },
    { name: 'blue', color: '#2196f3' },
    { name: 'lightblue', color: '#03a9f4' },
    { name: 'cyan', color: '#00bcd4' },
    { name: 'teal', color: '#009688' },
    { name: 'green', color: '#4caf50' },
    { name: 'lightgreen', color: '#8bc34a' },
    { name: 'lime', color: '#cddc39' },
    { name: 'yellow', color: '#ffeb3b' },
    { name: 'amber', color: '#ffc107' },
    { name: 'orange', color: '#fb8c00' },
    { name: 'deeporange', color: '#ff5722' },
    { name: 'brown', color: '#795548' },
    { name: 'grey', color: '#9e9e9e' },
    { name: 'bluegrey', color: '#607d8b' },
    { name: 'white', color: '#ffffff' },
  ];

  public WARN: Array<{ name: string; color: string; }> = [
    { name: 'red', color: '#f44336' },
    { name: 'pink', color: '#e91e63' },
    { name: 'purple', color: '#9c27b0' },
    { name: 'deeppurple', color: '#673ab7' },
  ];

  public FONTLIGHT: Array<{ name: string; color: string; }> = [
    { name: 'gray-500', color: '#6b7280' },
    { name: 'gray-600', color: '#4b5563' },
    { name: 'gray-700', color: '#374151' },
    { name: 'gray-800', color: '#1f2937' },
    { name: 'gray-900', color: '#111827' },
    { name: 'black', color: '#000000' },
  ];

  public FONTDARK: Array<{ name: string; color: string; }> = [
    { name: 'white', color: '#ffffff' },
    { name: 'gray-50', color: '#f9fafb' },
    { name: 'gray-100', color: '#f3f4f6' },
    { name: 'gray-200', color: '#e5e7eb' },
    { name: 'gray-300', color: '#d1d5db' },
    { name: 'gray-400', color: '#9ca3af' },
    { name: 'gray-500', color: '#6b7280' },
  ];

  public BACKGROUNDLIGHT: Array<{ name: string; color: string; }> = [
    { name: 'white', color: '#ffffff' },
    { name: 'gray-50', color: '#f9fafb' },
    { name: 'gray-100', color: '#f3f4f6' },
    { name: 'gray-200', color: '#e5e7eb' },
    { name: 'gray-300', color: '#d1d5db' },
    { name: 'gray-400', color: '#9ca3af' },
    { name: 'gray-500', color: '#6b7280' },
  ];

  public BACKGROUNDDARK: Array<{ name: string; color: string; }> = [
    { name: 'gray-500', color: '#6b7280' },
    { name: 'gray-600', color: '#4b5563' },
    { name: 'gray-700', color: '#374151' },
    { name: 'gray-800', color: '#1f2937' },
    { name: 'gray-900', color: '#111827' },
    { name: 'black', color: '#000000' },
  ];

  public colors: Array<{ name: string; color: string; }> = this.PRIMARY;
  public isOpen: boolean = false;

  @Input() colorType: ColorType = ColorType.PRIMARY;
  @Input() color: string = '';
  @Input() previewColor: string = '';
  @Input() name: string = '';
  @Output() previewChanged: EventEmitter<string> = new EventEmitter();

  public emitPreview(color: string): void {
    this.previewColor = color;
    this.previewChanged.emit(this.previewColor);
  }

  public ngOnInit(): void {
    switch (this.colorType) {
      case ColorType.PRIMARY:
        this.colors = this.PRIMARY;
        break;
      case ColorType.WARN:
        this.colors = this.WARN;
        break;
      case ColorType.FONTDARK:
        this.colors = this.FONTDARK;
        break;
      case ColorType.FONTLIGHT:
        this.colors = this.FONTLIGHT;
        break;
      case ColorType.BACKGROUNDDARK:
        this.colors = this.BACKGROUNDDARK;
        break;
      case ColorType.BACKGROUNDLIGHT:
        this.colors = this.BACKGROUNDLIGHT;
        break;
      default:
        this.colors = this.PRIMARY;
        break;
    }
  }

  public changeComplete(event: ColorEvent): void {
    this.emitPreview(event.color.hex);
  }

  public get previewColorCropped(): string {
    let s = this.previewColor;
    while (s.charAt(0) === '#') {
      s = s.substring(1);
    }
    return s;
  }
}
