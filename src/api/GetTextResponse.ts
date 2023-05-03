import MiResponse from "./MiResponse";
import Text from "./data_struct/Text";

export default class GetTextResponse extends MiResponse {
    public text: Text = new Text()
}