import MiResponse from "./MiResponse";

export default class GetTextsRelatedTaskResponse extends MiResponse {
    public texts: Array<Text> = new Array<Text>()
}