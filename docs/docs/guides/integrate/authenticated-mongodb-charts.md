---
title: Authenticated MongoDB Charts
---

This integration guide shows how you can embed authenticated MongoDB Charts in your web application using ZITADEL as authentication provider.

## Setup ZITADEL Application

Before you can embed an authenticated chart in your application, you have to do a few configuration steps in ZITADEL Console.
You will need to provide some information about your app. We recommend creating a new app to start from scratch.

1. Navigate to your Project
2. Add a new application at the top of the page.
3. Select Web application type and continue.
4. Use [Authorization Code](/apis/openidoauth/grant-types#authorization-code) in combination with [Proof Key for Code Exchange (PKCE)](/apis/openidoauth/grant-types#proof-key-for-code-exchange).
5. Skip the redirect settings and confirm the app creation
6. Copy the client ID, you will need to tell MongoDB Charts about it.
7. When you created the app, expand its _OIDC Configuration_ section, change the _Auth Token Type_ to _JWT_ and save the change.

Your application configuration should now look similar to this:

![Create app in console](/img/integrations/mongodb-charts-app-create-light.png)

## Setup Custom JWT Provider for MongoDB Charts

Configure ZITADEL as your _Custom JWT Provider_ following the [MongoDB docs](https://docs.mongodb.com/charts/configure-auth-providers/) .

Configure the following values:
- Signing Algorithm: RS256
- Signing Key: JWK or JWKS URL
- JWKS: https://{your_domain}.zitadel.cloud/oauth/v2/keys
- Audience: Your app's client ID which you copied when you created the ZITADEL app

Your configuration should look similar to this:

![Configure Custom JWT Provider](/img/integrations/mongodb-charts-auth-provider-light.png)

## Embedding your Chart

Embed a chart into your application now, following the corresponding [MongoDB docs](https://docs.mongodb.com/charts/saas/embed-chart-jwt-auth/).

If you've done the [Angular Quickstart](/examples/login/angular.md), your code could look something like this:

```html
<!-- chart.component.html -->
<div id="chart"></div>
```

```css
/* chart.component.css */
div#chart {
    height: 500px;    
}
```

```ts
// chart.component.ts
import { Component, OnInit } from '@angular/core';
import ChartsEmbedSDK from "@mongodb-js/charts-embed-dom";
import { AuthenticationService } from 'src/app/services/authentication.service';

@Component({
  selector: 'app-chart',
  templateUrl: './chart.component.html',
  styleUrls: ['./chart.component.css']
})
export class ChartComponent implements OnInit {

  constructor(private auth: AuthenticationService) { }

  ngOnInit(): void {
    this.renderChart().catch(e => window.alert(e.message));    
  }

  async renderChart() {
    const sdk = new ChartsEmbedSDK({
      baseUrl: "<YOUR CHARTS BASE URL HERE>",
      getUserToken: () => {
        return this.auth.getAccessToken()
      },
    });
  
    const chart = sdk.createChart({
      chartId: "<YOUR CHART ID HERE>"
    });
    await chart.render(<HTMLElement>document.getElementById("chart"));
  }  
}
```
