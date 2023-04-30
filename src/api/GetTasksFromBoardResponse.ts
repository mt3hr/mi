import MiResponse from "./MiResponse";
import TaskInfo from "./data_struct/TaskInfo";

export default class GetTasksFromBoardResponse extends MiResponse {
    public boards_tasks: Array<TaskInfo> = new Array<TaskInfo>()
}