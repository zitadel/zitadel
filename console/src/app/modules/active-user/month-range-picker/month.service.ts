import { computed, inject, Injectable, Signal } from "@angular/core";
import { TranslateService } from "@ngx-translate/core";
import { toSignal } from "@angular/core/rxjs-interop";
import { map } from "rxjs/operators";

@Injectable()
export class MonthService {
  public readonly monthNames: Signal<string[]>;

  constructor() {
    const translateService = inject(TranslateService);

    const lang = toSignal(
      translateService.onLangChange.pipe(map(({ lang }) => lang)),
      {
        initialValue: translateService.getCurrentLang(),
      }
    );

    this.monthNames = computed(() => {
      return Array.from({ length: 12 }, (_, i) =>
        new Intl.DateTimeFormat(lang(), { month: "short" }).format(
          new Date(2000, i)
        )
      );
    });
  }
}
