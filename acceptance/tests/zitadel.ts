import fetch from "node-fetch";

export async function removeUserByUsername(username: string) {
  const resp = await getUserByUsername(username);
  if (!resp || !resp.result || !resp.result[0]) {
    return;
  }
  await removeUser(resp.result[0].userId);
}

export async function removeUser(id: string) {
  const response = await fetch(process.env.ZITADEL_API_URL! + "/v2/users/" + id, {
    method: "DELETE",
    headers: {
      Authorization: "Bearer " + process.env.ZITADEL_SERVICE_USER_TOKEN!,
    },
  });
  if (response.statusCode >= 400 && response.statusCode != 404) {
    const error = "HTTP Error: " + response.statusCode + " - " + response.statusMessage;
    console.error(error);
    throw new Error(error);
  }
  return;
}

export async function getUserByUsername(username: string) {
  const listUsersBody = {
    queries: [
      {
        userNameQuery: {
          userName: username,
        },
      },
    ],
  };
  const jsonBody = JSON.stringify(listUsersBody);
  const registerResponse = await fetch(process.env.ZITADEL_API_URL! + "/v2/users", {
    method: "POST",
    body: jsonBody,
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + process.env.ZITADEL_SERVICE_USER_TOKEN!,
    },
  });
  if (registerResponse.statusCode >= 400) {
    const error = "HTTP Error: " + registerResponse.statusCode + " - " + registerResponse.statusMessage;
    console.error(error);
    throw new Error(error);
  }
  const respJson = await registerResponse.json();
  return respJson;
}
