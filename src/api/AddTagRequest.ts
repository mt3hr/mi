import MiRequest from "./MiRequest";

export default class AddTagRequest extends MiRequest {
    public task_id: string = ""
    public tag: string = ""
}