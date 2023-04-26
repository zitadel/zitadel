import { Injectable } from '@angular/core';
import { Meta } from '@angular/platform-browser';

@Injectable({
  providedIn: 'root',
})
export class SeoService {
  constructor(private meta: Meta) {}

  public generateTags(config: any): void {
    // default values
    config = {
      title: 'ZITADEL Console',
      description: 'Managementplatform for ZITADEL',
      image: 'https://www.zitadel.com/images/preview.png',
      slug: '',
      ...config,
    };

    this.meta.updateTag({ property: 'og:type', content: 'website' });
    this.meta.updateTag({ property: 'og:site_name', content: 'ZITADEL Console' });
    this.meta.updateTag({ property: 'og:title', content: config.title });
    this.meta.updateTag({ property: 'description', content: config.description });
    this.meta.updateTag({ property: 'og:description', content: config.description });
    if (config.image) {
      this.meta.updateTag({ property: 'og:image', content: config.image });
    }

    this.meta.updateTag({ property: 'twitter:card', content: 'summary' });
    this.meta.updateTag({ property: 'og:site', content: '@zitadel_ch' });
    this.meta.updateTag({ property: 'og:title', content: config.title });
    this.meta.updateTag({ property: 'og:image', content: 'https://www.zitadel.com/images/preview.png' });
    this.meta.updateTag({ property: 'og:description', content: config.description });
  }
}
