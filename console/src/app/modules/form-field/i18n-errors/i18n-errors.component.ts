import { KeyValue, KeyValuePipe } from '@angular/common';
import { Component, Input, OnInit } from '@angular/core';
import { AbstractControl } from '@angular/forms';
import { TranslateService } from '@ngx-translate/core';
import { filter, map, mergeMap, Observable } from 'rxjs';

@Component({
  selector: 'cnsl-i18n-errors',
  templateUrl: './i18n-errors.component.html',
})
export class I18nErrorsComponent implements OnInit {
  @Input() ctrl!: AbstractControl;

  public errors$!: Observable<Array<KeyValue<unknown, unknown>>>;

  constructor(
    private translateSvc: TranslateService,
    private kvPipe: KeyValuePipe,
    ) {}

    ngOnInit(): void {
      this.errors$ = this.ctrl.valueChanges.pipe(
        mergeMap(() => this.ctrl.statusChanges),
        filter(status => status === 'INVALID'),
        map(() => this.currentErrors())
      )
    }

  public currentErrors(): Array<KeyValue<unknown, unknown>> {
    return this.kvPipe.transform(this.ctrl.errors)?.filter(kv => {
      return kv.key as string != "invalid" && !(kv.value as any).invalid
    }) || []
  }

  public translate(err: KeyValue<unknown, unknown>): string {
    const anyValue = err.value as any
    return this.translateSvc.instant(anyValue.i18nKey || "ERRORS.INVALID_FORMAT")
  }
}
