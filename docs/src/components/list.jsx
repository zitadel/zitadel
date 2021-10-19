import React from 'react';

import styles from '../css/list.module.css';

export const ICONTYPE = {
  START: <div className="rounded rounded-start">
    <i className={`las la-play-circle`}></i>
  </div>,
  LOGIN: <div className="rounded rounded-login">
    <i className={`las la-sign-in-alt`}></i>
  </div>,
  PRIVATELABELING: <div className="rounded rounded-privatelabel">
    <i className={`las la-swatchbook`}></i>
  </div>,
  TEXTS: <div className="rounded rounded-texts">
    <i className={`las la-paragraph`}></i>
  </div>,
  SYSTEM: <div className="rounded rounded-system">
    <i className={`las la-server`}></i>
  </div>,
  APIS: <div className="rounded rounded-apis">
  <i className={`las la-code`}></i>
</div>
};

export function ListElement({ link, iconClasses, type, title, description}) {
  return (
    <a className={styles.listelement} href={link}>
      {type ? type : 
        iconClasses && <div><i className={`${styles.icon} ${iconClasses}`}></i></div>
      }
      <div>
        <h3>{title}</h3>
        <p className={styles.listelement.description}>{description}</p>
      </div>
    </a>
  )
}

export function ListWrapper({children, title}) {
  return (
    <div className={styles.listWrapper}>
      <span className={styles.listWrapperTitle}>{title}</span>
      {children}
    </div>
  )
}