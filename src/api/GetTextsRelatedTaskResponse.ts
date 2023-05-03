import MiResponse from "./MiResponse";
import Text from "./data_struct/Text";

export default class GetTextsRelatedTaskResponse extends MiResponse {
    public texts: Array<Text> = new Array<Text>()
}