import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { UntypedFormControl } from '@angular/forms';
import { Observable, Subscription } from 'rxjs';

@Component({
  selector: 'cnsl-additional-origins',
  templateUrl: './additional-origins.component.html',
  styleUrls: ['./additional-origins.component.scss'],
  standalone: false,
})
export class AdditionalOriginsComponent implements OnInit, OnDestroy {
  @Input() title: string = '';
  @Input() canWrite: boolean = false;
  @Input() public urisList: string[] = [];
  @Input() public redirectControl: UntypedFormControl = new UntypedFormControl({ value: '', disabled: true });
  @Output() public changedUris: EventEmitter<string[]> = new EventEmitter<string[]>();
  @Input() public getValues: Observable<void> = new Observable();
  public placeholder: string = '<scheme> "://" <hostname> [ ":" <port> ]';

  @ViewChild('originInput') input!: any;
  private sub: Subscription = new Subscription();

  constructor() {}

  ngOnInit(): void {
    if (this.canWrite) {
      this.redirectControl.enable();
    }

    this.sub = this.getValues.subscribe(() => {
      this.add(this.input.nativeElement);
    });
  }

  ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  public add(input: any): void {
    if (this.redirectControl.valid) {
      if (input.value !== '' && input.value !== ' ' && input.value !== '/') {
        this.urisList.push(input.value);
      }
      if (input) {
        input.value = '';
      }
    }
  }

  public remove(redirect: any): void {
    const index = this.urisList.indexOf(redirect);

    if (index >= 0) {
      this.urisList.splice(index, 1);
    }
  }
}
