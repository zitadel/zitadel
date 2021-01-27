import '../static/base.css';

import * as sapper from '@sapper/app';

import { startClient } from './i18n.js';
import { initPhotoSwipeFromDOM } from './utils/photoswipe.js';

startClient();

initPhotoSwipeFromDOM('.zitadel-gallery');

sapper.start({
    target: document.querySelector('#sapper')
});