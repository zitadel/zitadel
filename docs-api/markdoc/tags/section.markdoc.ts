import { Section } from "../../components/Section";

export default {
  render: Section,
  children: ["paragraph", "tag", "list", "code"],
  attributes: {
    columns: {
      type: Number,
      default: 1,
    },
  },
};
