import { Injectable } from '@angular/core';
import { Meta } from '@angular/platform-browser';

import { environment } from '../../environments/environment';

@Injectable({
    providedIn: 'root',
})
export class SeoService {
    constructor(private meta: Meta) { }

    public generateTags(config: any): void {
        // default values
        config = {
            title: 'ZITADEL Console',
            description: 'Managementplatform for ZITADEL',
            image: 'https://www.zitadel.ch/zitadel-social-preview25.png',
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
        this.meta.updateTag({ property: 'og:url', content: `https://${environment.production ? 'console.zitadel.ch' : 'console.zitadel.dev'}/${config.slug}` });

        this.meta.updateTag({ property: 'twitter:card', content: 'summary' });
        this.meta.updateTag({ property: 'og:site', content: '@zitadel_ch' });
        this.meta.updateTag({ property: 'og:title', content: config.title });
        this.meta.updateTag({ property: 'og:image', content: 'https://www.zitadel.ch/zitadel-social-preview25.png' });
        this.meta.updateTag({ property: 'og:description', content: config.description });
    }
}
