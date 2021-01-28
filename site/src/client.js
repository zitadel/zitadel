import '../static/base.css';

import * as sapper from '@sapper/app';

import { initPhotoSwipeFromDOM } from './utils/photoswipe.js';

initPhotoSwipeFromDOM('.zitadel-gallery');

sapper.start({
    target: document.querySelector('#sapper')
});