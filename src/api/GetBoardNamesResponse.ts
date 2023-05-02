import MiResponse from "./MiResponse";

export default class GetBoardNamesResponse extends MiResponse {
    public board_names: Array<string> = new Array<string>()
}