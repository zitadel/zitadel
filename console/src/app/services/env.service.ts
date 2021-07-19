import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class EnvService {
  private env$!: Observable<any>;
  constructor(private http: HttpClient) { }

  public loadEnvironment(): Observable<any> {
    if (!this.env$) {
      return this.http.get('./assets/environment.json');
    } else {
      console.log('loaded env from cache');
      return this.env$;
    }
  }
}
