import React from 'react';

import styles from '../css/apicard.module.css';

export function ApiCard({ title, type, label, children}) {
  let style = styles.apiauth;
  switch (type) {
    case 'AUTH':
      style = styles.apiauth;
      break;
    case 'MGMT':
      style = styles.apimgmt;
      break;
    case 'ADMIN':
      style = styles.apiadmin;
      break;
    case 'SYSTEM':
      style = styles.apisystem;
      break;
    case 'ASSET':
      style = styles.apiasset;
      break;      
  }
    
  return (
    <div className={`${styles.apicard} ${style} `}>
      {/* {title && <h2 className={styles.apicard.title}>{title}</h2>} */}
      {/* <p className={styles.apicard.description}>
        
      </p> */}
      {children}
    </div>
  )
}
