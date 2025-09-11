import { Location } from '@angular/common';
import { Injectable } from '@angular/core';
import { NavigationEnd, Router } from '@angular/router';

@Injectable({
  providedIn: 'root',
})
export class NavigationService {
  private history: string[] = [];

  constructor(
    private router: Router,
    private location: Location,
  ) {
    this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        this.history.push(event.urlAfterRedirects);
      }
    });
  }

  public back(): void {
    if (this.isBackPossible) {
      this.location.back();
      this.history.pop();
    } else {
      this.router.navigateByUrl('/');
      this.history.pop();
    }
  }

  public get isBackPossible(): boolean {
    return this.history.length > 1;
  }
}
