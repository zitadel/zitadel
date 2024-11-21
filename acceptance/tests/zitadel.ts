import axios from "axios";

export async function removeUserByUsername(username: string) {
  const resp = await getUserByUsername(username);
  if (!resp || !resp.result || !resp.result[0]) {
    return;
  }
  await removeUser(resp.result[0].userId);
}

export async function removeUser(id: string) {
  try {
    const response = await axios.delete(`${process.env.ZITADEL_API_URL}/v2/users/${id}`, {
      headers: {
        Authorization: `Bearer ${process.env.ZITADEL_SERVICE_USER_TOKEN}`,
      },
    });

    if (response.status >= 400 && response.status !== 404) {
      const error = `HTTP Error: ${response.status} - ${response.statusText}`;
      console.error(error);
      throw new Error(error);
    }
  } catch (error) {
    console.error("Error making request:", error);
    throw error;
  }
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

  try {
    const response = await axios.post(`${process.env.ZITADEL_API_URL}/v2/users`, listUsersBody, {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${process.env.ZITADEL_SERVICE_USER_TOKEN}`,
      },
    });

    if (response.status >= 400) {
      const error = `HTTP Error: ${response.status} - ${response.statusText}`;
      console.error(error);
      throw new Error(error);
    }

    return response.data;
  } catch (error) {
    console.error("Error making request:", error);
    throw error;
  }
}
