import MiResponse from "./MiResponse";
import Tag from "./data_struct/Tag";

export default class GetTagsRelatedTaskResponse extends MiResponse {
    public tags: Array<Tag> = new Array<Tag>()
}