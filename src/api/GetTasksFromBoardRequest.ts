import MiRequest from "./MiRequest";
import TaskSearchQuery from "./data_struct/TaskSearchQuery";

export default class GetTasksFromBoardRequest extends MiRequest {
    public query: TaskSearchQuery = new TaskSearchQuery()
    public update_cache: boolean = false
}