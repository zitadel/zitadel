import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

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
    { name: 'red', color: '#f44336' },
    { name: 'pink', color: '#e91e63' },
    { name: 'purple', color: '#9c27b0' },
    { name: 'deeppurple', color: '#673ab7' },
  ];

  public FONTDARK: Array<{ name: string; color: string; }> = [
    { name: 'red', color: '#f44336' },
    { name: 'pink', color: '#e91e63' },
    { name: 'purple', color: '#9c27b0' },
    { name: 'deeppurple', color: '#673ab7' },
  ];

  public BACKGROUNDLIGHT: Array<{ name: string; color: string; }> = [
    { name: 'red', color: '#f44336' },
    { name: 'pink', color: '#e91e63' },
    { name: 'purple', color: '#9c27b0' },
    { name: 'deeppurple', color: '#673ab7' },
  ];

  public BACKGROUNDDARK: Array<{ name: string; color: string; }> = [
    { name: 'red', color: '#f44336' },
    { name: 'pink', color: '#e91e63' },
    { name: 'purple', color: '#9c27b0' },
    { name: 'deeppurple', color: '#673ab7' },
  ];

  public colors = this.PRIMARY;

  @Input() colorType: ColorType = ColorType.PRIMARY;
  @Input() color: string = '';
  @Input() previewColor: string = '';
  @Input() name: string = '';
  @Output() previewChanged: EventEmitter<string> = new EventEmitter();

  public emitPreview(color: string): void {
    this.previewColor = color;
    this.previewChanged.emit(this.previewColor);
  }

  ngOnInit(): void {

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
}
