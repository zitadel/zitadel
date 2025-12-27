import { Component, EventEmitter, Input, Output, ViewChild, ElementRef } from '@angular/core';

@Component({
  selector: 'cnsl-top-view',
  templateUrl: './top-view.component.html',
  styleUrls: ['./top-view.component.scss'],
  standalone: false,
})
export class TopViewComponent {
  @Input() public title: string = '';
  @Input() public sub: string = '';
  @Input() public stateTooltip: string = '';
  @Input() public isActive: boolean = false;
  @Input() public isInactive: boolean = false;
  @Input() public hasActions: boolean | null = false;
  @Input() public hasContributors: boolean | null = false;
  @Input() public docLink: string = '';
  @Input() public hasBackButton: boolean | null = true;
  @Input() public titleEditable: boolean = false;
  @Input() public titlePlaceholder: string = '';
  @Output() public backClicked: EventEmitter<void> = new EventEmitter();
  @Output() public titleChanged: EventEmitter<string> = new EventEmitter();

  @ViewChild('titleEditInput') set titleInput(element: ElementRef<HTMLInputElement>) {
    if (element) {
      element.nativeElement.focus();
    }
  }

  public isEditingTitle: boolean = false;
  public editingTitleValue: string = '';

  constructor() {}

  public backClick(): void {
    this.backClicked.emit();
  }

  public startTitleEdit(): void {
    if (this.titleEditable) {
      this.isEditingTitle = true;
      this.editingTitleValue = this.title;
      // Focus the input after the view updates
      setTimeout(() => {
        const input = document.querySelector('.cnsl-title-input') as HTMLInputElement;
        if (input) {
          input.focus();
          input.select();
        }
      });
    }
  }

  public saveTitleEdit(): void {
    if (this.editingTitleValue.trim() && this.editingTitleValue !== this.title) {
      this.titleChanged.emit(this.editingTitleValue.trim());
    }
    this.cancelTitleEdit();
  }

  public cancelTitleEdit(): void {
    this.isEditingTitle = false;
    this.editingTitleValue = '';
  }
}
