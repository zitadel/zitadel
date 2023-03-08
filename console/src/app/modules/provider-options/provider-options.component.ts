import { Component, EventEmitter, Input, OnChanges, OnDestroy, Output } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { Subject, takeUntil } from 'rxjs';
import { Options } from 'src/app/proto/generated/zitadel/idp_pb';

@Component({
  selector: 'cnsl-provider-options',
  templateUrl: './provider-options.component.html',
  styleUrls: ['./provider-options.component.scss'],
})
export class ProviderOptionsComponent implements OnChanges, OnDestroy {
  @Input() public initialOptions?: Options.AsObject;
  @Output() public optionsChanged: EventEmitter<Options> = new EventEmitter<Options>();
  private destroy$: Subject<void> = new Subject();
  public form: FormGroup = new FormGroup({
    isAutoCreation: new FormControl(false, []),
    isAutoUpdate: new FormControl(false, []),
    isCreationAllowed: new FormControl(true, []),
    isLinkingAllowed: new FormControl(true, []),
  });

  constructor() {
    this.form.valueChanges.pipe(takeUntil(this.destroy$)).subscribe((value) => {
      if (value) {
        const opt = new Options();
        opt.setIsAutoCreation(value.isAutoCreation);
        opt.setIsAutoUpdate(value.isAutoUpdate);
        opt.setIsCreationAllowed(value.isCreationAllowed);
        opt.setIsLinkingAllowed(value.isLinkingAllowed);
        this.optionsChanged.emit(opt);
      }
    });
  }

  ngOnChanges(): void {
    if (this.initialOptions) {
      this.form.patchValue(this.initialOptions);
    }
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
