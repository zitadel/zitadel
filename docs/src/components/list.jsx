import React from "react";

import styles from "../css/list.module.css";

export const ICONTYPE = {
  START: (
    <div className="custom-rounded custom-rounded-start">
      <i className={`las la-play-circle`}></i>
    </div>
  ),
  TASKS: (
    <div className="custom-rounded custom-rounded-start">
      <i className={`las la-tasks`}></i>
    </div>
  ),
  ARCHITECTURE: (
    <div className="custom-rounded custom-rounded-architecture">
      <i className={`las la-sitemap`}></i>
    </div>
  ),
  INSTANCE: (
    <div className="custom-rounded custom-rounded-instance">
      <i className={`las la-industry`}></i>
    </div>
  ),
  LOGIN: (
    <div className="custom-rounded custom-rounded-login">
      <i className={`las la-sign-in-alt`}></i>
    </div>
  ),
  PRIVATELABELING: (
    <div className="custom-rounded custom-rounded-privatelabel">
      <i className={`las la-swatchbook`}></i>
    </div>
  ),
  TEXTS: (
    <div className="custom-rounded custom-rounded-texts">
      <i className={`las la-paragraph`}></i>
    </div>
  ),
  POLICY: (
    <div className="custom-rounded custom-rounded-policy">
      <i className={`las la-file-contract`}></i>
    </div>
  ),
  SERVICE: (
    <div className="custom-rounded custom-rounded-service">
      <i className={`las la-concierge-bell`}></i>
    </div>
  ),
  STORAGE: (
    <div className="custom-rounded custom-rounded-storage">
      <i className={`las la-database`}></i>
    </div>
  ),
  FOLDER: (
    <div className="custom-rounded custom-rounded-storage">
      <i className={`las la-folder`}></i>
    </div>
  ),
  FILE: (
    <div className="custom-rounded custom-rounded-storage">
      <i className={`las la-file-alt`}></i>
    </div>
  ),
  SYSTEM: (
    <div className="custom-rounded custom-rounded-system">
      <i className={`las la-server`}></i>
    </div>
  ),
  APIS: (
    <div className="custom-rounded custom-rounded-apis">
      <i className={`las la-code`}></i>
    </div>
  ),
  HELP: (
    <div className="custom-rounded custom-rounded-help">
      <i className={`las la-question`}></i>
    </div>
  ),
  HELP_REGISTER: (
    <div className="custom-rounded custom-rounded-login">
      <i className={`las la-plus-circle`}></i>
    </div>
  ),
  HELP_LOGIN: (
    <div className="custom-rounded custom-rounded-login">
      <i className={`las la-sign-in-alt`}></i>
    </div>
  ),
  HELP_PASSWORDLESS: (
    <div className="custom-rounded custom-rounded-login">
      <i className={`las la-fingerprint`}></i>
    </div>
  ),
  HELP_PASSWORD: (
    <div className="custom-rounded custom-rounded-password">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        xmlnsXlink="http://www.w3.org/1999/xlink"
        version="1.1"
        width="100%"
        height="100%"
        viewBox="0 0 24 24"
        fit=""
        preserveAspectRatio="xMidYMid meet"
        focusable="false"
      >
        <path d="M17,7H22V17H17V19A1,1 0 0,0 18,20H20V22H17.5C16.95,22 16,21.55 16,21C16,21.55 15.05,22 14.5,22H12V20H14A1,1 0 0,0 15,19V5A1,1 0 0,0 14,4H12V2H14.5C15.05,2 16,2.45 16,3C16,2.45 16.95,2 17.5,2H20V4H18A1,1 0 0,0 17,5V7M2,7H13V9H4V15H13V17H2V7M20,15V9H17V15H20M8.5,12A1.5,1.5 0 0,0 7,10.5A1.5,1.5 0 0,0 5.5,12A1.5,1.5 0 0,0 7,13.5A1.5,1.5 0 0,0 8.5,12M13,10.89C12.39,10.33 11.44,10.38 10.88,11C10.32,11.6 10.37,12.55 11,13.11C11.55,13.63 12.43,13.63 13,13.11V10.89Z"></path>
      </svg>
    </div>
  ),
  HELP_FACTORS: (
    <div className="custom-rounded custom-rounded-login">
      <i className={`las la-fingerprint`}></i>
    </div>
  ),
  HELP_PHONE: (
    <div className="custom-rounded custom-rounded-phone">
      <i className={`las la-phone`}></i>
    </div>
  ),
  HELP_EMAIL: (
    <div className="custom-rounded custom-rounded-email">
      <i className={`las la-at`}></i>
    </div>
  ),
  HELP_SOCIAL: (
    <div className="custom-rounded custom-rounded-login">
      <i className={`las la-share-alt`}></i>
    </div>
  ),
};

export function ListElement({
  link,
  iconClasses,
  roundClasses,
  label,
  type,
  title,
  description,
}) {
  return (
    <a 
      className={styles.listelement} 
      href={link} 
      onClick={() => {
        window.plausible("ListElement", {
          props: { method: title },
        });
      }}
    >
      {type
        ? type
        : iconClasses && (
            <div className={roundClasses}>
              {label ? (
                <span className={styles.listlabel}>{label}</span>
              ) : (
                <i className={`${iconClasses}`}></i>
              )}
            </div>
          )}
      <div>
        <h3>{title}</h3>
        <p className={styles.listelement.description}>{description}</p>
      </div>
    </a>
  );
}

export function ListWrapper({ children, title, columns }) {
  return (
    <div className={styles.listWrapper}>
      {title && <span className={styles.listWrapperTitle}>{title}</span>}
      {children}
    </div>
  );
}

export function HomeListWrapper({ children, image }) {
  return (
    <div className={styles.homerow}>
      {image}
      <div className={styles.homecontent}>{children}</div>
    </div>
  );
}
