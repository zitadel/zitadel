import React from 'react';

import styles from '../css/card.module.css';

export function Card({ link, githubLink, imageSource, title, description, label}) {
  if (githubLink) {
    return (
      <a className={styles.card} href={githubLink} target="_blank">
        <h3 className={styles.card.title}>{title}</h3>
        {description && <p className={styles.card.description}>{description}</p>}
        <span className={styles.fillspace}></span>
        <div className={styles.bottom}>
          <img className={styles.bottomicon} src="/docs/img/tech/github.svg" alt="github"/>
          <span className={styles.bottomspan}>{label}</span>
        </div>
      </a>
    )
  } else if (link){
    return (
      <a className={styles.card} href={link}>
        {imageSource && <img src={imageSource} className={styles.cardimg} alt={`image ${title}`}/>}
        <h3 className={styles.card.title}>{title}</h3>
        <p className={styles.card.description}>{description}</p>
      </a>
    )
  };
}

export function CardWrapper({children}) {
  return (
    <div className={styles.cardWrapper}>
      {children}
    </div>
  )
}