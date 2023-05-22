"use server";

import { cookies } from "next/headers";

export type SessionCookie = {
  id: string;
  token: string;
  loginName: string;
  changeDate: string;
};

async function set(sessions: SessionCookie[]) {
  const cookiesList = cookies();
  // @ts-ignore
  cookiesList.set({
    name: "sessions",
    value: JSON.stringify(sessions),
    httpOnly: true,
    path: "/",
  });
}

export async function addSessionToCookie(session: SessionCookie): Promise<any> {
  const cookiesList = cookies();
  //   const hasSessions = cookiesList.has("sessions");
  //   if (hasSessions) {
  const stringifiedCookie = cookiesList.get("sessions");

  const currentSessions: SessionCookie[] = stringifiedCookie?.value
    ? JSON.parse(stringifiedCookie?.value)
    : [];

  // @ts-ignore
  return cookiesList.set({
    name: "sessions",
    value: JSON.stringify([...currentSessions, session]),
    httpOnly: true,
    path: "/",
  });
  //   } else {
  //     return set([session]);
  //   }
}

export async function updateSessionCookie(
  id: string,
  session: SessionCookie
): Promise<any> {
  const cookiesList = cookies();
  //   const hasSessions = cookiesList.has("sessions");
  //   if (hasSessions) {
  const stringifiedCookie = cookiesList.get("sessions");

  const sessions: SessionCookie[] = stringifiedCookie?.value
    ? JSON.parse(stringifiedCookie?.value)
    : [session];

  const foundIndex = sessions.findIndex((session) => session.id === id);
  sessions[foundIndex] = session;

  // @ts-ignore
  return cookiesList.set({
    name: "sessions",
    value: JSON.stringify(sessions),
    httpOnly: true,
    path: "/",
  });
  //   } else {
  //     return Promise.reject();
  //   }
}

export async function removeSessionFromCookie(
  session: SessionCookie
): Promise<any> {
  const cookiesList = cookies();
  //   const hasSessions = cookiesList.has("sessions");
  //   if (hasSessions) {
  const stringifiedCookie = cookiesList.get("sessions");

  const sessions: SessionCookie[] = stringifiedCookie?.value
    ? JSON.parse(stringifiedCookie?.value)
    : [session];

  const filteredSessions = sessions.filter(
    (session) => session.id !== session.id
  );

  // @ts-ignore
  return cookiesList.set({
    name: "__Secure-sessions",
    value: JSON.stringify(filteredSessions),
    httpOnly: true,
    path: "/",
  });
  //   } else {
  //     return Promise.reject();
  //   }
}

export async function getMostRecentSessionCookie(): Promise<any> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie[] = JSON.parse(stringifiedCookie?.value);

    console.log(sessions);
    const latest = sessions.reduce((prev, current) => {
      return new Date(prev.changeDate).getTime() >
        new Date(current.changeDate).getTime()
        ? prev
        : current;
    });

    return latest;
  } else {
    return Promise.reject();
  }
}

export async function getAllSessionIds(): Promise<any> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie[] = JSON.parse(stringifiedCookie?.value);
    return sessions.map((session) => session.id);
  } else {
    return Promise.reject();
  }
}

/**
 * Returns most recent session filtered by optinal loginName
 * @param loginName
 * @returns most recent session
 */
export async function getMostRecentCookieWithLoginname(
  loginName?: string
): Promise<any> {
  const cookiesList = cookies();

  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie[] = JSON.parse(stringifiedCookie?.value);

    const filtered = sessions.filter((cookie) => {
      console.log(!!loginName);
      return !!loginName ? cookie.loginName === loginName : true;
    });

    console.log(filtered);

    const latest =
      filtered && filtered.length
        ? filtered.reduce((prev, current) => {
            return new Date(prev.changeDate).getTime() >
              new Date(current.changeDate).getTime()
              ? prev
              : current;
          })
        : undefined;

    if (latest) {
      return latest;
    } else {
      return Promise.reject();
    }
  } else {
    return Promise.reject();
  }
}

export async function clearSessions() {}
