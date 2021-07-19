import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class EnvService {
  env: any;
  constructor(private http: HttpClient) {
    this.loadEnvironment();
  }

  public loadEnvironment(): Promise<any> {
    if (this.env) {
      console.log('loaded env from cache');
      return Promise.resolve(this.env);
    } else {
      return this.http.get('./assets/environment.json')
        .toPromise().then((data: any) => {
          this.env = data;
          return this.env;
        });
    }
  }
}
