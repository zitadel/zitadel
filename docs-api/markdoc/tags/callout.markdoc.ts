import { Callout } from "../../components/Callout";

export default {
  render: Callout,
  children: ["paragraph", "tag", "list"],
  attributes: {
    title: {
      type: String,
    },
  },
};
