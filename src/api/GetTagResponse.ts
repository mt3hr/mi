import MiResponse from "./MiResponse";
import Tag from "./data_struct/Tag";

export default class GetTagResponse extends MiResponse {
    public tag: Tag = new Tag()
}