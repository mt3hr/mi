import MiResponse from "./MiResponse";

export default class GetTagNamesResponse extends MiResponse {
    public tag_names: Array<string> = new Array<string>()
}