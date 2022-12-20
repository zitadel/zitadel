import { Endpoint } from "../../components/Endpoint";

export default {
  render: Endpoint,
  children: ["paragraph"],
  attributes: {
    method: {
      type: String,
      default: "GET",
    },
    link: {
      type: String,
    },
  },
};
