import { Injectable } from '@angular/core';
import { Event, NavigationEnd, Router } from '@angular/router';
import { filter } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class RefreshService {
  private fifo: Array<string> = ['', ''];
  constructor(router: Router) {
    router.events.pipe(filter((event) => event instanceof NavigationEnd)).subscribe((event: Event | any) => {
      this.moveInto(event.url);
    });
  }

  public get previousUrls(): Array<string> {
    return this.fifo;
  }

  private moveInto(value: string): void {
    this.fifo[1] = this.fifo[0];
    this.fifo[0] = value;
  }
}
